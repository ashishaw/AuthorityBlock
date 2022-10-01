// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package builtin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ashishaw/authorityblock/abi"
	"github.com/ashishaw/authorityblock/ablock"
	"github.com/ashishaw/authorityblock/xenv"
)

func init() {

	events := Prototype.Events()

	mustEventByName := func(name string) *abi.Event {
		if event, found := events.EventByName(name); found {
			return event
		}
		panic("event not found")
	}

	masterEvent := mustEventByName("$Master")
	creditPlanEvent := mustEventByName("$CreditPlan")
	userEvent := mustEventByName("$User")
	sponsorEvent := mustEventByName("$Sponsor")

	defines := []struct {
		name string
		run  func(env *xenv.Environment) []interface{}
	}{
		{"native_master", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)

			env.UseGas(ablock.GetBalanceGas)
			master, err := env.State().GetMaster(ablock.Address(self))
			if err != nil {
				panic(err)
			}

			return []interface{}{master}
		}},
		{"native_setMaster", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self      common.Address
				NewMaster common.Address
			}
			env.ParseArgs(&args)

			env.UseGas(ablock.SstoreResetGas)
			if err := env.State().SetMaster(ablock.Address(args.Self), ablock.Address(args.NewMaster)); err != nil {
				panic(err)
			}

			env.Log(masterEvent, ablock.Address(args.Self), nil, args.NewMaster)
			return nil
		}},
		{"native_balanceAtBlock", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self        common.Address
				BlockNumber uint32
			}
			env.ParseArgs(&args)
			ctx := env.BlockContext()

			if args.BlockNumber > ctx.Number {
				return []interface{}{&big.Int{}}
			}

			if ctx.Number-args.BlockNumber > ablock.MaxStateHistory {
				return []interface{}{&big.Int{}}
			}

			if args.BlockNumber == ctx.Number {
				env.UseGas(ablock.GetBalanceGas)
				val, err := env.State().GetBalance(ablock.Address(args.Self))
				if err != nil {
					panic(err)
				}
				return []interface{}{val}
			}

			env.UseGas(ablock.SloadGas)
			env.UseGas(ablock.SloadGas)
			summary, err := env.Chain().GetBlockSummary(args.BlockNumber)
			if err != nil {
				panic(err)
			}

			env.UseGas(ablock.SloadGas)
			state := env.State().Checkout(summary.Header.StateRoot(), summary.Header.Number(), summary.Conflicts, summary.SteadyNum)

			env.UseGas(ablock.GetBalanceGas)
			val, err := state.GetBalance(ablock.Address(args.Self))
			if err != nil {
				panic(err)
			}

			return []interface{}{val}
		}},
		{"native_energyAtBlock", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self        common.Address
				BlockNumber uint32
			}
			env.ParseArgs(&args)
			ctx := env.BlockContext()
			if args.BlockNumber > ctx.Number {
				return []interface{}{&big.Int{}}
			}

			if ctx.Number-args.BlockNumber > ablock.MaxStateHistory {
				return []interface{}{&big.Int{}}
			}

			if args.BlockNumber == ctx.Number {
				env.UseGas(ablock.GetBalanceGas)
				val, err := env.State().GetEnergy(ablock.Address(args.Self), ctx.Time)
				if err != nil {
					panic(err)
				}
				return []interface{}{val}
			}

			env.UseGas(ablock.SloadGas)
			env.UseGas(ablock.SloadGas)
			summary, err := env.Chain().GetBlockSummary(args.BlockNumber)
			if err != nil {
				panic(err)
			}

			env.UseGas(ablock.SloadGas)
			state := env.State().Checkout(summary.Header.StateRoot(), summary.Header.Number(), summary.Conflicts, summary.SteadyNum)

			env.UseGas(ablock.GetBalanceGas)
			val, err := state.GetEnergy(ablock.Address(args.Self), summary.Header.Timestamp())
			if err != nil {
				panic(err)
			}

			return []interface{}{val}
		}},
		{"native_hasCode", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)

			env.UseGas(ablock.GetBalanceGas)
			codeHash, err := env.State().GetCodeHash(ablock.Address(self))
			if err != nil {
				panic(err)
			}
			hasCode := !codeHash.IsZero()

			return []interface{}{hasCode}
		}},
		{"native_storageFor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				Key  ablock.Bytes32
			}
			env.ParseArgs(&args)

			env.UseGas(ablock.SloadGas)
			storage, err := env.State().GetStorage(ablock.Address(args.Self), args.Key)
			if err != nil {
				panic(err)
			}
			return []interface{}{storage}
		}},
		{"native_creditPlan", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(self))

			env.UseGas(ablock.SloadGas)
			credit, rate, err := binding.CreditPlan()
			if err != nil {
				panic(err)
			}

			return []interface{}{credit, rate}
		}},
		{"native_setCreditPlan", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self         common.Address
				Credit       *big.Int
				RecoveryRate *big.Int
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SstoreSetGas)
			if err := binding.SetCreditPlan(args.Credit, args.RecoveryRate); err != nil {
				panic(err)
			}
			env.Log(creditPlanEvent, ablock.Address(args.Self), nil, args.Credit, args.RecoveryRate)
			return nil
		}},
		{"native_isUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isUser, err := binding.IsUser(ablock.Address(args.User))
			if err != nil {
				panic(err)
			}

			return []interface{}{isUser}
		}},
		{"native_userCredit", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(2 * ablock.SloadGas)
			credit, err := binding.UserCredit(ablock.Address(args.User), env.BlockContext().Time)
			if err != nil {
				panic(err)
			}

			return []interface{}{credit}
		}},
		{"native_addUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isUser, err := binding.IsUser(ablock.Address(args.User))
			if err != nil {
				panic(err)
			}
			if isUser {
				return []interface{}{false}
			}

			env.UseGas(ablock.SstoreSetGas)
			if err := binding.AddUser(ablock.Address(args.User), env.BlockContext().Time); err != nil {
				panic(err)
			}

			var action ablock.Bytes32
			copy(action[:], "added")
			env.Log(userEvent, ablock.Address(args.Self), []ablock.Bytes32{ablock.BytesToBytes32(args.User[:])}, action)
			return []interface{}{true}
		}},
		{"native_removeUser", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self common.Address
				User common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isUser, err := binding.IsUser(ablock.Address(args.User))
			if err != nil {
				panic(err)
			}
			if !isUser {
				return []interface{}{false}
			}

			env.UseGas(ablock.SstoreResetGas)
			if err := binding.RemoveUser(ablock.Address(args.User)); err != nil {
				panic(err)
			}

			var action ablock.Bytes32
			copy(action[:], "removed")
			env.Log(userEvent, ablock.Address(args.Self), []ablock.Bytes32{ablock.BytesToBytes32(args.User[:])}, action)
			return []interface{}{true}
		}},
		{"native_sponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isSponsor, err := binding.IsSponsor(ablock.Address(args.Sponsor))
			if err != nil {
				panic(err)
			}
			if isSponsor {
				return []interface{}{false}
			}

			env.UseGas(ablock.SstoreSetGas)
			if err := binding.Sponsor(ablock.Address(args.Sponsor), true); err != nil {
				panic(err)
			}

			var action ablock.Bytes32
			copy(action[:], "sponsored")
			env.Log(sponsorEvent, ablock.Address(args.Self), []ablock.Bytes32{ablock.BytesToBytes32(args.Sponsor.Bytes())}, action)
			return []interface{}{true}
		}},
		{"native_unsponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isSponsor, err := binding.IsSponsor(ablock.Address(args.Sponsor))
			if err != nil {
				panic(err)
			}
			if !isSponsor {
				return []interface{}{false}
			}

			env.UseGas(ablock.SstoreResetGas)
			if err := binding.Sponsor(ablock.Address(args.Sponsor), false); err != nil {
				panic(err)
			}

			var action ablock.Bytes32
			copy(action[:], "unsponsored")
			env.Log(sponsorEvent, ablock.Address(args.Self), []ablock.Bytes32{ablock.BytesToBytes32(args.Sponsor.Bytes())}, action)
			return []interface{}{true}
		}},
		{"native_isSponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isSponsor, err := binding.IsSponsor(ablock.Address(args.Sponsor))
			if err != nil {
				panic(err)
			}

			return []interface{}{isSponsor}
		}},
		{"native_selectSponsor", func(env *xenv.Environment) []interface{} {
			var args struct {
				Self    common.Address
				Sponsor common.Address
			}
			env.ParseArgs(&args)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(args.Self))

			env.UseGas(ablock.SloadGas)
			isSponsor, err := binding.IsSponsor(ablock.Address(args.Sponsor))
			if err != nil {
				panic(err)
			}
			if !isSponsor {
				return []interface{}{false}
			}

			env.UseGas(ablock.SstoreResetGas)
			binding.SelectSponsor(ablock.Address(args.Sponsor))

			var action ablock.Bytes32
			copy(action[:], "selected")
			env.Log(sponsorEvent, ablock.Address(args.Self), []ablock.Bytes32{ablock.BytesToBytes32(args.Sponsor.Bytes())}, action)

			return []interface{}{true}
		}},
		{"native_currentSponsor", func(env *xenv.Environment) []interface{} {
			var self common.Address
			env.ParseArgs(&self)
			binding := Prototype.Native(env.State()).Bind(ablock.Address(self))

			env.UseGas(ablock.SloadGas)
			addr, err := binding.CurrentSponsor()
			if err != nil {
				panic(err)
			}

			return []interface{}{addr}
		}},
	}
	abi := Prototype.NativeABI()
	for _, def := range defines {
		if method, found := abi.MethodByName(def.name); found {
			nativeMethods[methodKey{Prototype.Address, method.ID()}] = &nativeMethod{
				abi: method,
				run: def.run,
			}
		} else {
			panic("method not found: " + def.name)
		}
	}
}
