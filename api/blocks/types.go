// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package blocks

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/tx"
)

type BFTEngine interface {
	Finalized() ablock.Bytes32
}

type JSONBlockSummary struct {
	Number       uint32       `json:"number"`
	ID           ablock.Bytes32 `json:"id"`
	Size         uint32       `json:"size"`
	ParentID     ablock.Bytes32 `json:"parentID"`
	Timestamp    uint64       `json:"timestamp"`
	GasLimit     uint64       `json:"gasLimit"`
	Beneficiary  ablock.Address `json:"beneficiary"`
	GasUsed      uint64       `json:"gasUsed"`
	TotalScore   uint64       `json:"totalScore"`
	TxsRoot      ablock.Bytes32 `json:"txsRoot"`
	TxsFeatures  uint32       `json:"txsFeatures"`
	StateRoot    ablock.Bytes32 `json:"stateRoot"`
	ReceiptsRoot ablock.Bytes32 `json:"receiptsRoot"`
	COM          bool         `json:"com"`
	Signer       ablock.Address `json:"signer"`
	IsTrunk      bool         `json:"isTrunk"`
	IsFinalized  bool         `json:"isFinalized"`
}

type JSONCollapsedBlock struct {
	*JSONBlockSummary
	Transactions []ablock.Bytes32 `json:"transactions"`
}

type JSONClause struct {
	To    *ablock.Address        `json:"to"`
	Value math.HexOrDecimal256 `json:"value"`
	Data  string               `json:"data"`
}

type JSONTransfer struct {
	Sender    ablock.Address          `json:"sender"`
	Recipient ablock.Address          `json:"recipient"`
	Amount    *math.HexOrDecimal256 `json:"amount"`
}

type JSONEvent struct {
	Address ablock.Address   `json:"address"`
	Topics  []ablock.Bytes32 `json:"topics"`
	Data    string         `json:"data"`
}

type JSONOutput struct {
	ContractAddress *ablock.Address   `json:"contractAddress"`
	Events          []*JSONEvent    `json:"events"`
	Transfers       []*JSONTransfer `json:"transfers"`
}

type JSONEmbeddedTx struct {
	ID           ablock.Bytes32        `json:"id"`
	ChainTag     byte                `json:"chainTag"`
	BlockRef     string              `json:"blockRef"`
	Expiration   uint32              `json:"expiration"`
	Clauses      []*JSONClause       `json:"clauses"`
	GasPriceCoef uint8               `json:"gasPriceCoef"`
	Gas          uint64              `json:"gas"`
	Origin       ablock.Address        `json:"origin"`
	Delegator    *ablock.Address       `json:"delegator"`
	Nonce        math.HexOrDecimal64 `json:"nonce"`
	DependsOn    *ablock.Bytes32       `json:"dependsOn"`
	Size         uint32              `json:"size"`

	// receipt part
	GasUsed  uint64                `json:"gasUsed"`
	GasPayer ablock.Address          `json:"gasPayer"`
	Paid     *math.HexOrDecimal256 `json:"paid"`
	Reward   *math.HexOrDecimal256 `json:"reward"`
	Reverted bool                  `json:"reverted"`
	Outputs  []*JSONOutput         `json:"outputs"`
}

type JSONExpandedBlock struct {
	*JSONBlockSummary
	Transactions []*JSONEmbeddedTx `json:"transactions"`
}

func buildJSONBlockSummary(summary *chain.BlockSummary, isTrunk bool, isFinalized bool) *JSONBlockSummary {
	header := summary.Header
	signer, _ := header.Signer()

	return &JSONBlockSummary{
		Number:       header.Number(),
		ID:           header.ID(),
		ParentID:     header.ParentID(),
		Timestamp:    header.Timestamp(),
		TotalScore:   header.TotalScore(),
		GasLimit:     header.GasLimit(),
		GasUsed:      header.GasUsed(),
		Beneficiary:  header.Beneficiary(),
		Signer:       signer,
		Size:         uint32(summary.Size),
		StateRoot:    header.StateRoot(),
		ReceiptsRoot: header.ReceiptsRoot(),
		TxsRoot:      header.TxsRoot(),
		TxsFeatures:  uint32(header.TxsFeatures()),
		COM:          header.COM(),
		IsTrunk:      isTrunk,
		IsFinalized:  isFinalized,
	}
}

func buildJSONOutput(txID ablock.Bytes32, index uint32, c *tx.Clause, o *tx.Output) *JSONOutput {
	jo := &JSONOutput{
		ContractAddress: nil,
		Events:          make([]*JSONEvent, 0, len(o.Events)),
		Transfers:       make([]*JSONTransfer, 0, len(o.Transfers)),
	}
	if c.To() == nil {
		addr := ablock.CreateContractAddress(txID, index, 0)
		jo.ContractAddress = &addr
	}
	for _, e := range o.Events {
		jo.Events = append(jo.Events, &JSONEvent{
			Address: e.Address,
			Data:    hexutil.Encode(e.Data),
			Topics:  e.Topics,
		})
	}
	for _, t := range o.Transfers {
		jo.Transfers = append(jo.Transfers, &JSONTransfer{
			Sender:    t.Sender,
			Recipient: t.Recipient,
			Amount:    (*math.HexOrDecimal256)(t.Amount),
		})
	}
	return jo
}

func buildJSONEmbeddedTxs(txs tx.Transactions, receipts tx.Receipts) []*JSONEmbeddedTx {
	jTxs := make([]*JSONEmbeddedTx, 0, len(txs))
	for itx, tx := range txs {
		receipt := receipts[itx]

		clauses := tx.Clauses()
		blockRef := tx.BlockRef()
		origin, _ := tx.Origin()
		delegator, _ := tx.Delegator()

		jcs := make([]*JSONClause, 0, len(clauses))
		jos := make([]*JSONOutput, 0, len(receipt.Outputs))

		for i, c := range clauses {
			jcs = append(jcs, &JSONClause{
				c.To(),
				math.HexOrDecimal256(*c.Value()),
				hexutil.Encode(c.Data()),
			})
			if !receipt.Reverted {
				jos = append(jos, buildJSONOutput(tx.ID(), uint32(i), c, receipt.Outputs[i]))
			}
		}

		jTxs = append(jTxs, &JSONEmbeddedTx{
			ID:           tx.ID(),
			ChainTag:     tx.ChainTag(),
			BlockRef:     hexutil.Encode(blockRef[:]),
			Expiration:   tx.Expiration(),
			Clauses:      jcs,
			GasPriceCoef: tx.GasPriceCoef(),
			Gas:          tx.Gas(),
			Origin:       origin,
			Delegator:    delegator,
			Nonce:        math.HexOrDecimal64(tx.Nonce()),
			DependsOn:    tx.DependsOn(),
			Size:         uint32(tx.Size()),

			GasUsed:  receipt.GasUsed,
			GasPayer: receipt.GasPayer,
			Paid:     (*math.HexOrDecimal256)(receipt.Paid),
			Reward:   (*math.HexOrDecimal256)(receipt.Reward),
			Reverted: receipt.Reverted,
			Outputs:  jos,
		})
	}
	return jTxs
}
