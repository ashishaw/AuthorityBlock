// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package node

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ashishaw/authorityblock/ablock"
)

type Master struct {
	PrivateKey  *ecdsa.PrivateKey
	Beneficiary *ablock.Address
}

func (m *Master) Address() ablock.Address {
	return ablock.Address(crypto.PubkeyToAddress(m.PrivateKey.PublicKey))
}
