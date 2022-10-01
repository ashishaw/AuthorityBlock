// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package poa

import (
	"github.com/ashishaw/authorityblock/block"
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/ablock"
)

var epochInterval uint32 = ablock.SeederInterval

// mockEpochInterval mocks the epoch interval。
// TEST ONLY
func mockEpochInterval(interval uint32) {
	epochInterval = interval
}

// Seeder generates seed for poa scheduler.
type Seeder struct {
	repo  *chain.Repository
	cache map[ablock.Bytes32][]byte
}

// NewSeeder creates a seeder
func NewSeeder(repo *chain.Repository) *Seeder {
	return &Seeder{
		repo,
		make(map[ablock.Bytes32][]byte),
	}
}

// Generate creates a seed for the given parent block's header.
func (seeder *Seeder) Generate(parentID ablock.Bytes32) (seed []byte, err error) {
	blockNum := block.Number(parentID) + 1

	epoch := blockNum / epochInterval
	if epoch <= 1 {
		return
	}
	seedNum := (epoch - 1) * epochInterval

	seedID, err := seeder.repo.NewChain(parentID).GetBlockID(seedNum)
	if err != nil {
		return
	}

	if v, ok := seeder.cache[seedID]; ok {
		return v, nil
	}

	defer func() {
		if err == nil {
			seeder.cache[seedID] = seed
		}
	}()

	summary, err := seeder.repo.GetBlockSummary(seedID)
	if err != nil {
		return
	}

	return summary.Header.Beta()
}
