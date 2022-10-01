// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package transfers

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ashishaw/authorityblock/api/events"
	"github.com/ashishaw/authorityblock/logdb"
	"github.com/ashishaw/authorityblock/ablock"
)

type LogMeta struct {
	BlockID        ablock.Bytes32 `json:"blockID"`
	BlockNumber    uint32       `json:"blockNumber"`
	BlockTimestamp uint64       `json:"blockTimestamp"`
	TxID           ablock.Bytes32 `json:"txID"`
	TxOrigin       ablock.Address `json:"txOrigin"`
	ClauseIndex    uint32       `json:"clauseIndex"`
}

type FilteredTransfer struct {
	Sender    ablock.Address          `json:"sender"`
	Recipient ablock.Address          `json:"recipient"`
	Amount    *math.HexOrDecimal256 `json:"amount"`
	Meta      LogMeta               `json:"meta"`
}

func convertTransfer(transfer *logdb.Transfer) *FilteredTransfer {
	v := math.HexOrDecimal256(*transfer.Amount)
	return &FilteredTransfer{
		Sender:    transfer.Sender,
		Recipient: transfer.Recipient,
		Amount:    &v,
		Meta: LogMeta{
			BlockID:        transfer.BlockID,
			BlockNumber:    transfer.BlockNumber,
			BlockTimestamp: transfer.BlockTime,
			TxID:           transfer.TxID,
			TxOrigin:       transfer.TxOrigin,
			ClauseIndex:    transfer.ClauseIndex,
		},
	}
}

type TransferFilter struct {
	CriteriaSet []*logdb.TransferCriteria
	Range       *events.Range
	Options     *logdb.Options
	Order       logdb.Order //default asc
}
