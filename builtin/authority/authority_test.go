// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package authority

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/ablock"
)

func M(a ...interface{}) []interface{} {
	return a
}

func TestAuthority(t *testing.T) {
	db := muxdb.NewMem()
	st := state.New(db, ablock.Bytes32{}, 0, 0, 0)

	p1 := ablock.BytesToAddress([]byte("p1"))
	p2 := ablock.BytesToAddress([]byte("p2"))
	p3 := ablock.BytesToAddress([]byte("p3"))

	st.SetBalance(p1, big.NewInt(10))
	st.SetBalance(p2, big.NewInt(20))
	st.SetBalance(p3, big.NewInt(30))

	aut := New(ablock.BytesToAddress([]byte("aut")), st)
	tests := []struct {
		ret      interface{}
		expected interface{}
	}{
		{M(aut.Add(p1, p1, ablock.Bytes32{})), M(true, nil)},
		{M(aut.Get(p1)), M(true, p1, ablock.Bytes32{}, true, nil)},
		{M(aut.Add(p2, p2, ablock.Bytes32{})), M(true, nil)},
		{M(aut.Add(p3, p3, ablock.Bytes32{})), M(true, nil)},
		{M(aut.Candidates(big.NewInt(10), ablock.InitialMaxBlockProposers)), M(
			[]*Candidate{{p1, p1, ablock.Bytes32{}, true}, {p2, p2, ablock.Bytes32{}, true}, {p3, p3, ablock.Bytes32{}, true}}, nil,
		)},
		{M(aut.Candidates(big.NewInt(20), ablock.InitialMaxBlockProposers)), M(
			[]*Candidate{{p2, p2, ablock.Bytes32{}, true}, {p3, p3, ablock.Bytes32{}, true}}, nil,
		)},
		{M(aut.Candidates(big.NewInt(30), ablock.InitialMaxBlockProposers)), M(
			[]*Candidate{{p3, p3, ablock.Bytes32{}, true}}, nil,
		)},
		{M(aut.Candidates(big.NewInt(10), 2)), M(
			[]*Candidate{{p1, p1, ablock.Bytes32{}, true}, {p2, p2, ablock.Bytes32{}, true}}, nil,
		)},
		{M(aut.Get(p1)), M(true, p1, ablock.Bytes32{}, true, nil)},
		{M(aut.Update(p1, false)), M(true, nil)},
		{M(aut.Get(p1)), M(true, p1, ablock.Bytes32{}, false, nil)},
		{M(aut.Update(p1, true)), M(true, nil)},
		{M(aut.Get(p1)), M(true, p1, ablock.Bytes32{}, true, nil)},
		{M(aut.Revoke(p1)), M(true, nil)},
		{M(aut.Get(p1)), M(false, p1, ablock.Bytes32{}, false, nil)},
		{M(aut.Candidates(&big.Int{}, ablock.InitialMaxBlockProposers)), M(
			[]*Candidate{{p2, p2, ablock.Bytes32{}, true}, {p3, p3, ablock.Bytes32{}, true}}, nil,
		)},
	}

	for i, tt := range tests {
		assert.Equal(t, tt.expected, tt.ret, "#%v", i)
	}
}
