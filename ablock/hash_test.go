// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package ablock_test

import (
	"hash"
	"io"
	"math/rand"
	"testing"

	"github.com/ashishaw/authorityblock/ablock"
	"golang.org/x/crypto/sha3"
)

func BenchmarkHash(b *testing.B) {
	data := make([]byte, 10)
	rand.New(rand.NewSource(1)).Read(data)

	b.Run("keccak", func(b *testing.B) {
		type keccakState interface {
			hash.Hash
			Read([]byte) (int, error)
		}

		k := sha3.NewLegacyKeccak256().(keccakState)
		var b32 ablock.Bytes32
		for i := 0; i < b.N; i++ {
			k.Write(data)
			k.Read(b32[:])
			k.Reset()
		}
	})

	b.Run("blake2b", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ablock.Blake2b(data)
		}
	})
}

func BenchmarkBlake2b(b *testing.B) {
	data := make([]byte, 100)
	rand.New(rand.NewSource(1)).Read(data)
	b.Run("Blake2b", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ablock.Blake2b(data).Bytes()
		}
	})

	b.Run("BlakeFn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ablock.Blake2bFn(func(w io.Writer) {
				w.Write(data)
			})
		}
	})
}
