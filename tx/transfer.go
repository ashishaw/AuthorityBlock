// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package tx

import (
	"math/big"

	"github.com/ashishaw/authorityblock/ablock"
)

// Transfer token transfer log.
type Transfer struct {
	Sender    ablock.Address
	Recipient ablock.Address
	Amount    *big.Int
}

// Transfers slisce of transfer logs.
type Transfers []*Transfer
