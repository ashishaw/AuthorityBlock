// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package solo

import (
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/comm"
	"github.com/ashishaw/authorityblock/ablock"
)

// Communicator in solo is a fake one just for api handler.
type Communicator struct {
}

// PeersStats returns nil solo doesn't join p2p network.
func (comm *Communicator) PeersStats() []*comm.PeerStats {
	return nil
}

// BFTEngine is a fake bft engine for solo.
type BFTEngine struct {
	finalized ablock.Bytes32
}

func (engine *BFTEngine) Finalized() ablock.Bytes32 {
	return engine.finalized
}

func NewBFTEngine(repo *chain.Repository) *BFTEngine {
	return &BFTEngine{
		finalized: repo.GenesisBlock().Header().ID(),
	}
}
