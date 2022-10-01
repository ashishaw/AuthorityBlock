// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package state

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/ablock"
)

func TestCachedObject(t *testing.T) {
	db := muxdb.NewMem()
	addr := ablock.Address{}

	stgTrie := db.NewTrie(StorageTrieName([]byte("sid")), ablock.Bytes32{}, 0, 0)
	storages := []struct {
		k ablock.Bytes32
		v rlp.RawValue
	}{
		{ablock.BytesToBytes32([]byte("key1")), []byte("value1")},
		{ablock.BytesToBytes32([]byte("key2")), []byte("value2")},
		{ablock.BytesToBytes32([]byte("key3")), []byte("value3")},
		{ablock.BytesToBytes32([]byte("key4")), []byte("value4")},
	}

	for _, s := range storages {
		saveStorage(stgTrie, s.k, s.v)
	}

	storageRoot, commit := stgTrie.Stage(0, 0)

	err := commit()
	assert.Nil(t, err)

	code := make([]byte, 100)
	rand.Read(code)

	codeHash := crypto.Keccak256(code)
	db.NewStore(codeStoreName).Put(codeHash, code)

	account := Account{
		Balance:     &big.Int{},
		CodeHash:    codeHash,
		StorageRoot: storageRoot[:],
	}

	obj := newCachedObject(db, addr, &account, &AccountMetadata{StorageID: []byte("sid")})

	assert.Equal(t,
		M(code, nil),
		M(obj.GetCode()))

	for _, s := range storages {
		assert.Equal(t,
			M(s.v, nil),
			M(obj.GetStorage(s.k, 0)))
	}
}
