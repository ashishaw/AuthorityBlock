// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package poa

import (
	"github.com/ashishaw/authorityblock/ablock"
)

// Proposer address with status.
type Proposer struct {
	Address ablock.Address
	Active  bool
}
