// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/ashishaw/authorityblock/abi"
	"github.com/ashishaw/authorityblock/builtin/gen"
	"github.com/ashishaw/authorityblock/ablock"
)

type contract struct {
	name    string
	Address ablock.Address
	ABI     *abi.ABI
}

func mustLoadContract(name string) *contract {
	asset := "compiled/" + name + ".abi"
	data := gen.MustAsset(asset)
	abi, err := abi.New(data)
	if err != nil {
		panic(errors.Wrap(err, "load ABI for '"+name+"'"))
	}

	return &contract{
		name,
		ablock.BytesToAddress([]byte(name)),
		abi,
	}
}

// RuntimeBytecodes load runtime byte codes.
func (c *contract) RuntimeBytecodes() []byte {
	asset := "compiled/" + c.name + ".bin-runtime"
	data, err := hex.DecodeString(string(gen.MustAsset(asset)))
	if err != nil {
		panic(errors.Wrap(err, "load runtime byte code for '"+c.name+"'"))
	}
	return data
}

func (c *contract) NativeABI() *abi.ABI {
	asset := "compiled/" + c.name + "Native.abi"
	data := gen.MustAsset(asset)
	abi, err := abi.New(data)
	if err != nil {
		panic(errors.Wrap(err, "load native ABI for '"+c.name+"'"))
	}
	return abi
}
