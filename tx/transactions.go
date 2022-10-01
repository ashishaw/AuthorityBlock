// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package tx

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/trie"
)

var (
	emptyRoot = trie.DeriveRoot(&derivableTxs{})
)

// Transactions a slice of transactions.
type Transactions []*Transaction

// RootHash computes merkle root hash of transactions.
func (txs Transactions) RootHash() ablock.Bytes32 {
	if len(txs) == 0 {
		// optimized
		return emptyRoot
	}
	return trie.DeriveRoot(derivableTxs(txs))
}

// implements types.DerivableList
type derivableTxs Transactions

func (txs derivableTxs) Len() int {
	return len(txs)
}

func (txs derivableTxs) GetRlp(i int) []byte {
	data, err := rlp.EncodeToBytes(txs[i])
	if err != nil {
		panic(err)
	}
	return data
}
