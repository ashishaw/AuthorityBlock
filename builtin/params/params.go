// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/ablock"
)

// Params binder of `Params` contract.
type Params struct {
	addr  ablock.Address
	state *state.State
}

func New(addr ablock.Address, state *state.State) *Params {
	return &Params{addr, state}
}

// Get native way to get param.
func (p *Params) Get(key ablock.Bytes32) (value *big.Int, err error) {
	err = p.state.DecodeStorage(p.addr, key, func(raw []byte) error {
		if len(raw) == 0 {
			value = &big.Int{}
			return nil
		}
		return rlp.DecodeBytes(raw, &value)
	})
	return
}

// Set native way to set param.
func (p *Params) Set(key ablock.Bytes32, value *big.Int) error {
	return p.state.EncodeStorage(p.addr, key, func() ([]byte, error) {
		if value.Sign() == 0 {
			return nil, nil
		}
		return rlp.EncodeToBytes(value)
	})
}
