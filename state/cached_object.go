// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package state

import (
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/ablock"
)

var codeCache, _ = lru.NewARC(512)

// cachedObject to cache code and storage of an account.
type cachedObject struct {
	db   *muxdb.MuxDB
	addr ablock.Address
	data Account
	meta AccountMetadata

	cache struct {
		code        []byte
		storageTrie *muxdb.Trie
		storage     map[ablock.Bytes32]rlp.RawValue
	}
}

func newCachedObject(db *muxdb.MuxDB, addr ablock.Address, data *Account, meta *AccountMetadata) *cachedObject {
	return &cachedObject{db: db, addr: addr, data: *data, meta: *meta}
}

func (co *cachedObject) getOrCreateStorageTrie() *muxdb.Trie {
	if co.cache.storageTrie != nil {
		return co.cache.storageTrie
	}

	if len(co.data.StorageRoot) == 0 {
		return nil
	}

	trie := co.db.NewTrie(
		StorageTrieName(co.meta.StorageID),
		ablock.BytesToBytes32(co.data.StorageRoot),
		co.meta.StorageCommitNum,
		co.meta.StorageDistinctNum)

	co.cache.storageTrie = trie
	return trie
}

// GetStorage returns storage value for given key.
func (co *cachedObject) GetStorage(key ablock.Bytes32, steadyBlockNum uint32) (rlp.RawValue, error) {
	cache := &co.cache
	// retrive from storage cache
	if cache.storage != nil {
		if v, ok := cache.storage[key]; ok {
			return v, nil
		}
	} else {
		cache.storage = make(map[ablock.Bytes32]rlp.RawValue)
	}
	// not found in cache

	trie := co.getOrCreateStorageTrie()
	if trie == nil {
		return nil, nil
	}

	// load from trie
	v, err := loadStorage(trie, key, steadyBlockNum)
	if err != nil {
		return nil, err
	}
	// put into cache
	cache.storage[key] = v
	return v, nil
}

// GetCode returns the code of the account.
func (co *cachedObject) GetCode() ([]byte, error) {
	cache := &co.cache

	if len(cache.code) > 0 {
		return cache.code, nil
	}

	if len(co.data.CodeHash) > 0 {
		// do have code
		if code, has := codeCache.Get(string(co.data.CodeHash)); has {
			return code.([]byte), nil
		}

		code, err := co.db.NewStore(codeStoreName).Get(co.data.CodeHash)
		if err != nil {
			return nil, err
		}
		codeCache.Add(string(co.data.CodeHash), code)
		cache.code = code
		return code, nil
	}
	return nil, nil
}