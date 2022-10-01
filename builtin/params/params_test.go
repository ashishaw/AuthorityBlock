// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package params

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/ablock"
)

func TestParamsGetSet(t *testing.T) {
	db := muxdb.NewMem()
	st := state.New(db, ablock.Bytes32{}, 0, 0, 0)
	setv := big.NewInt(10)
	key := ablock.BytesToBytes32([]byte("key"))
	p := New(ablock.BytesToAddress([]byte("par")), st)
	p.Set(key, setv)

	getv, err := p.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, setv, getv)
}
