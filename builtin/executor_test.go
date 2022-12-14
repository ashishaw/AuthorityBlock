// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin_test

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/builtin"
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/runtime"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/tx"
	"github.com/ashishaw/authorityblock/xenv"
)

func M(a ...interface{}) []interface{} {
	return a
}

func approverEvent(approver ablock.Address, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("Approver")
	var b32 ablock.Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  []ablock.Bytes32{ev.ID(), ablock.BytesToBytes32(approver.Bytes())},
		Data:    data,
	}
}
func proposalEvent(id ablock.Bytes32, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("Proposal")
	var b32 ablock.Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  []ablock.Bytes32{ev.ID(), id},
		Data:    data,
	}
}

func votingContractEvent(addr ablock.Address, action string) *tx.Event {
	ev, _ := builtin.Executor.ABI.EventByName("VotingContract")
	var b32 ablock.Bytes32
	copy(b32[:], action)
	data, _ := ev.Encode(b32)
	return &tx.Event{
		Address: builtin.Executor.Address,
		Topics:  []ablock.Bytes32{ev.ID(), ablock.BytesToBytes32(addr.Bytes())},
		Data:    data,
	}
}

func initExectorTest() *ctest {
	db := muxdb.NewMem()
	b0 := buildGenesis(db, func(state *state.State) error {
		state.SetCode(builtin.Prototype.Address, builtin.Prototype.RuntimeBytecodes())
		state.SetCode(builtin.Executor.Address, builtin.Executor.RuntimeBytecodes())
		state.SetCode(builtin.Params.Address, builtin.Params.RuntimeBytecodes())
		builtin.Params.Native(state).Set(ablock.KeyExecutorAddress, new(big.Int).SetBytes(builtin.Executor.Address[:]))
		return nil
	})

	repo, _ := chain.NewRepository(db, b0)
	st := state.New(db, b0.Header().StateRoot(), 0, 0, 0)
	chain := repo.NewChain(b0.Header().ID())

	rt := runtime.New(chain, st, &xenv.BlockContext{Time: uint64(time.Now().Unix())}, ablock.NoFork)

	return &ctest{
		rt:  rt,
		abi: builtin.Executor.ABI,
		to:  builtin.Executor.Address,
	}
}

func TestExecutorApprover(t *testing.T) {
	test := initExectorTest()
	var approvers []ablock.Address
	for i := 0; i < 7; i++ {
		approvers = append(approvers, ablock.BytesToAddress([]byte(fmt.Sprintf("approver%d", i))))
	}

	for _, a := range approvers {
		// zero identity
		test.Case("addApprover", a, ablock.Bytes32{}).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("addApprover", a, ablock.BytesToBytes32(a.Bytes())).
			Caller(ablock.BytesToAddress([]byte("other"))).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("addApprover", a, ablock.BytesToBytes32(a.Bytes())).
			Caller(builtin.Executor.Address).
			ShouldLog(approverEvent(a, "added")).
			Assert(t)
		assert.Equal(t, M(true, nil), M(builtin.Prototype.Native(test.rt.State()).Bind(test.to).IsUser(a)))
	}

	test.Case("approverCount").
		ShouldOutput(uint8(len(approvers))).
		Assert(t)

	test.Case("addApprover", approvers[0], ablock.BytesToBytes32(approvers[0].Bytes())).
		Caller(builtin.Executor.Address).
		ShouldVMError(errReverted).
		Assert(t)

	for _, a := range approvers {
		test.Case("approvers", a).
			ShouldOutput(ablock.BytesToBytes32(a.Bytes()), true).
			Assert(t)
	}

	for _, a := range approvers {
		test.Case("revokeApprover", a).
			ShouldVMError(errReverted).
			Assert(t)

		test.Case("revokeApprover", a).
			Caller(builtin.Executor.Address).
			ShouldLog(approverEvent(a, "revoked")).
			Assert(t)
		assert.Equal(t, M(false, nil), M(builtin.Prototype.Native(test.rt.State()).Bind(test.to).IsUser(a)))
	}
	test.Case("approverCount").
		ShouldOutput(uint8(0)).
		Assert(t)
}

func TestExecutorVotingContract(t *testing.T) {

	test := initExectorTest()
	voting := ablock.BytesToAddress([]byte("voting"))
	test.Case("attachVotingContract", voting).
		ShouldVMError(errReverted).
		Assert(t)
	test.Case("votingContracts", voting).
		ShouldOutput(false).
		Assert(t)
	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "attached")).
		Assert(t)

	test.Case("votingContracts", voting).
		ShouldOutput(true).
		Assert(t)

	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("detachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "detached")).
		Assert(t)

	test.Case("attachVotingContract", voting).
		Caller(builtin.Executor.Address).
		ShouldLog(votingContractEvent(voting, "attached")).
		Assert(t)
}

func TestExecutorProposal(t *testing.T) {
	test := initExectorTest()

	target := builtin.Params.Address
	setParam, _ := builtin.Params.ABI.MethodByName("set")
	data, _ := setParam.EncodeInput(ablock.BytesToBytes32([]byte("paramKey")), big.NewInt(12345))
	test.Case("propose", target, data).
		ShouldVMError(errReverted).
		Assert(t)

	approver := ablock.BytesToAddress([]byte("approver"))
	test.Case("addApprover", approver, ablock.BytesToBytes32(approver.Bytes())).
		Caller(builtin.Executor.Address).
		Assert(t)

	proposalID := func() ablock.Bytes32 {
		var b8 [8]byte
		binary.BigEndian.PutUint64(b8[:], test.rt.Context().Time)
		return ablock.Bytes32(crypto.Keccak256Hash(b8[:], approver[:]))
	}()
	test.Case("propose", target, data).
		Caller(approver).
		ShouldOutput(proposalID).
		ShouldLog(proposalEvent(proposalID, "proposed")).
		Assert(t)

	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(0),
			false,
			target,
			data).
		Assert(t)

	test.Case("approve", proposalID).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldVMError(errReverted).
		Assert(t)

	test.Case("approve", proposalID).
		Caller(approver).
		ShouldLog(proposalEvent(proposalID, "approved")).
		Assert(t)
	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(1),
			false,
			target,
			data).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldLog(proposalEvent(proposalID, "executed")).
		Assert(t)

	test.Case("execute", proposalID).
		ShouldVMError(errReverted).
		Assert(t)
	test.Case("proposals", proposalID).
		ShouldOutput(
			test.rt.Context().Time,
			approver,
			uint8(1),
			uint8(1),
			true,
			target,
			data).
		Assert(t)

	assert.Equal(t, M(big.NewInt(12345), nil), M(builtin.Params.Native(test.rt.State()).Get(ablock.BytesToBytes32([]byte("paramKey")))))
}
