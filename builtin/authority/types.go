// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package authority

import (
	"github.com/ashishaw/authorityblock/ablock"
)

type (
	entry struct {
		Endorsor ablock.Address
		Identity ablock.Bytes32
		Active   bool
		Prev     *ablock.Address `rlp:"nil"`
		Next     *ablock.Address `rlp:"nil"`
	}

	// Candidate candidate of block proposer.
	Candidate struct {
		NodeMaster ablock.Address
		Endorsor   ablock.Address
		Identity   ablock.Bytes32
		Active     bool
	}
)

// IsEmpty returns whether the entry can be treated as empty.
func (e *entry) IsEmpty() bool {
	return e.Endorsor.IsZero() &&
		e.Identity.IsZero() &&
		!e.Active &&
		e.Prev == nil &&
		e.Next == nil
}

func (e *entry) IsLinked() bool {
	return e.Prev != nil || e.Next != nil
}
