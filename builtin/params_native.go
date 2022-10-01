// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/xenv"
)

func init() {
	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_executor", func(env *xenv.Environment) []interface{} {
			env.UseGas(ablock.SloadGas)
			val, err := Params.Native(env.State()).Get(ablock.KeyExecutorAddress)
			if err != nil {
				panic(err)
			}
			addr := ablock.BytesToAddress(val.Bytes())
			return []interface{}{addr}
		}},
		{"native_get", func(env *xenv.Environment) []interface{} {
			var key common.Hash
			env.ParseArgs(&key)

			env.UseGas(ablock.SloadGas)
			v, err := Params.Native(env.State()).Get(ablock.Bytes32(key))
			if err != nil {
				panic(err)
			}
			return []interface{}{v}
		}},
		{"native_set", func(env *xenv.Environment) []interface{} {
			var args struct {
				Key   common.Hash
				Value *big.Int
			}
			env.ParseArgs(&args)

			env.UseGas(ablock.SstoreSetGas)
			if err := Params.Native(env.State()).Set(ablock.Bytes32(args.Key), args.Value); err != nil {
				panic(err)
			}
			return nil
		}},
	}
	abi := Params.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Params.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
