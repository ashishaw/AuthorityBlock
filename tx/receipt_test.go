// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package tx_test

import (
	"fmt"
	"testing"

	. "github.com/ashishaw/authorityblock/tx"
)

func TestReceipt(t *testing.T) {
	var rs Receipts
	fmt.Println(rs.RootHash())

	var txs Transactions
	fmt.Println(txs.RootHash())
}
