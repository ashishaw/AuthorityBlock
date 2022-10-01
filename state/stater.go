// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package state

import (
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/ablock"
)

// Stater is the state creator.
type Stater struct {
	db *muxdb.MuxDB
}

// NewStater create a new stater.
func NewStater(db *muxdb.MuxDB) *Stater {
	return &Stater{db}
}

// NewState create a new state object.
func (s *Stater) NewState(root ablock.Bytes32, blockNum, blockConflicts, steadyBlockNum uint32) *State {
	return New(s.db, root, blockNum, blockConflicts, steadyBlockNum)
}
