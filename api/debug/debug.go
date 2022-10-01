// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package debug

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/ashishaw/authorityblock/api/utils"
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/consensus"
	"github.com/ashishaw/authorityblock/genesis"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/runtime"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/tracers"
	"github.com/ashishaw/authorityblock/tracers/logger"
	"github.com/ashishaw/authorityblock/trie"
	"github.com/ashishaw/authorityblock/vm"
)

var devNetGenesisID = genesis.NewDevnet().ID()

type Debug struct {
	repo       *chain.Repository
	stater     *state.Stater
	forkConfig ablock.ForkConfig
}

func New(repo *chain.Repository, stater *state.Stater, forkConfig ablock.ForkConfig) *Debug {
	return &Debug{
		repo,
		stater,
		forkConfig,
	}
}

func (d *Debug) handleTxEnv(ctx context.Context, blockID ablock.Bytes32, txIndex uint64, clauseIndex uint64) (*runtime.Runtime, *runtime.TransactionExecutor, error) {
	block, err := d.repo.GetBlock(blockID)
	if err != nil {
		if d.repo.IsNotFound(err) {
			return nil, nil, utils.Forbidden(errors.New("block not found"))
		}
		return nil, nil, err
	}
	txs := block.Transactions()
	if txIndex >= uint64(len(txs)) {
		return nil, nil, utils.Forbidden(errors.New("tx index out of range"))
	}
	if clauseIndex >= uint64(len(txs[txIndex].Clauses())) {
		return nil, nil, utils.Forbidden(errors.New("clause index out of range"))
	}
	skipPoA := d.repo.GenesisBlock().Header().ID() == devNetGenesisID
	rt, err := consensus.New(
		d.repo,
		d.stater,
		d.forkConfig,
	).NewRuntimeForReplay(block.Header(), skipPoA)
	if err != nil {
		return nil, nil, err
	}
	for i, tx := range txs {
		if uint64(i) > txIndex {
			break
		}
		txExec, err := rt.PrepareTransaction(tx)
		if err != nil {
			return nil, nil, err
		}
		clauseCounter := uint64(0)
		for txExec.HasNextClause() {
			if txIndex == uint64(i) && clauseIndex == clauseCounter {
				return rt, txExec, nil
			}
			if _, _, err := txExec.NextClause(); err != nil {
				return nil, nil, err
			}
			clauseCounter++
		}
		if _, err := txExec.Finalize(); err != nil {
			return nil, nil, err
		}
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}
	}
	return nil, nil, utils.Forbidden(errors.New("early reverted"))
}

//trace an existed transaction
func (d *Debug) traceTransaction(ctx context.Context, tracer tracers.Tracer, blockID ablock.Bytes32, txIndex uint64, clauseIndex uint64) (interface{}, error) {
	rt, txExec, err := d.handleTxEnv(ctx, blockID, txIndex, clauseIndex)
	if err != nil {
		return nil, err
	}
	rt.SetVMConfig(vm.Config{Debug: true, Tracer: tracer})
	_, _, err = txExec.NextClause()
	if err != nil {
		return nil, err
	}
	return tracer.GetResult()
}

func (d *Debug) handleTraceTransaction(w http.ResponseWriter, req *http.Request) error {
	var opt *TracerOption
	if err := utils.ParseJSON(req.Body, &opt); err != nil {
		return utils.BadRequest(errors.WithMessage(err, "body"))
	}
	if opt == nil {
		return utils.BadRequest(errors.New("body: empty body"))
	}
	var tracer tracers.Tracer
	if opt.Name == "" {
		tr, err := logger.NewStructLogger(opt.Config)
		if err != nil {
			return err
		}
		tracer = tr
	} else {
		name := opt.Name
		if !strings.HasSuffix(name, "Tracer") {
			name += "Tracer"
		}
		tr, err := tracers.New(name, nil, opt.Config)
		if err != nil {
			return err
		}
		tracer = tr
	}
	blockID, txIndex, clauseIndex, err := d.parseTarget(opt.Target)
	if err != nil {
		return err
	}
	res, err := d.traceTransaction(req.Context(), tracer, blockID, txIndex, clauseIndex)
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, res)
}

func (d *Debug) debugStorage(ctx context.Context, contractAddress ablock.Address, blockID ablock.Bytes32, txIndex uint64, clauseIndex uint64, keyStart []byte, maxResult int) (*StorageRangeResult, error) {
	rt, _, err := d.handleTxEnv(ctx, blockID, txIndex, clauseIndex)
	if err != nil {
		return nil, err
	}
	storageTrie, err := rt.State().BuildStorageTrie(contractAddress)
	if err != nil {
		return nil, err
	}
	return storageRangeAt(storageTrie, keyStart, maxResult)
}

func storageRangeAt(t *muxdb.Trie, start []byte, maxResult int) (*StorageRangeResult, error) {
	it := trie.NewIterator(t.NodeIterator(start, 0))
	result := StorageRangeResult{Storage: StorageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return nil, err
		}
		v := ablock.BytesToBytes32(content)
		e := StorageEntry{Value: &v}
		preimage := ablock.BytesToBytes32(it.Meta)
		e.Key = &preimage
		result.Storage[ablock.BytesToBytes32(it.Key).String()] = e
	}
	if it.Next() {
		next := ablock.BytesToBytes32(it.Key)
		result.NextKey = &next
	}
	return &result, nil
}

func (d *Debug) handleDebugStorage(w http.ResponseWriter, req *http.Request) error {
	var opt *StorageRangeOption
	if err := utils.ParseJSON(req.Body, &opt); err != nil {
		return utils.BadRequest(errors.WithMessage(err, "body"))
	}
	if opt == nil {
		return utils.BadRequest(errors.New("body: empty body"))
	}
	blockID, txIndex, clauseIndex, err := d.parseTarget(opt.Target)
	if err != nil {
		return err
	}
	var keyStart []byte
	if opt.KeyStart != "" {
		k, err := hexutil.Decode(opt.KeyStart)
		if err != nil {
			return utils.BadRequest(errors.New("keyStart: invalid format"))
		}
		keyStart = k
	}
	res, err := d.debugStorage(req.Context(), opt.Address, blockID, txIndex, clauseIndex, keyStart, opt.MaxResult)
	if err != nil {
		return err
	}
	return utils.WriteJSON(w, res)
}

func (d *Debug) parseTarget(target string) (blockID ablock.Bytes32, txIndex uint64, clauseIndex uint64, err error) {
	parts := strings.Split(target, "/")
	if len(parts) != 3 {
		return ablock.Bytes32{}, 0, 0, utils.BadRequest(errors.New("target:" + target + " unsupported"))
	}
	blockID, err = ablock.ParseBytes32(parts[0])
	if err != nil {
		return ablock.Bytes32{}, 0, 0, utils.BadRequest(errors.WithMessage(err, "target[0]"))
	}
	if len(parts[1]) == 64 || len(parts[1]) == 66 {
		txID, err := ablock.ParseBytes32(parts[1])
		if err != nil {
			return ablock.Bytes32{}, 0, 0, utils.BadRequest(errors.WithMessage(err, "target[1]"))
		}

		txMeta, err := d.repo.NewChain(blockID).GetTransactionMeta(txID)
		if err != nil {
			if d.repo.IsNotFound(err) {
				return ablock.Bytes32{}, 0, 0, utils.Forbidden(errors.New("transaction not found"))
			}
			return ablock.Bytes32{}, 0, 0, err
		}
		txIndex = txMeta.Index
	} else {
		i, err := strconv.ParseUint(parts[1], 0, 0)
		if err != nil {
			return ablock.Bytes32{}, 0, 0, utils.BadRequest(errors.WithMessage(err, "target[1]"))
		}
		txIndex = i
	}
	clauseIndex, err = strconv.ParseUint(parts[2], 0, 0)
	if err != nil {
		return ablock.Bytes32{}, 0, 0, utils.BadRequest(errors.WithMessage(err, "target[2]"))
	}
	return
}

func (d *Debug) Mount(root *mux.Router, pathPrefix string) {
	sub := root.PathPrefix(pathPrefix).Subrouter()

	sub.Path("/tracers").Methods(http.MethodPost).HandlerFunc(utils.WrapHandlerFunc(d.handleTraceTransaction))
	sub.Path("/storage-range").Methods(http.MethodPost).HandlerFunc(utils.WrapHandlerFunc(d.handleDebugStorage))

}
