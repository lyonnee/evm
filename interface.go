// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"math/big"

	"github.com/lyonnee/evm/common"
)

type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	// 添加退款金额
	AddRefund(uint64)
	// 减去退款金额
	SubRefund(uint64)
	GetRefund() uint64

	// 返回的是当前执行context提交到数据库的已提交状态(committed state)
	// 也就是上次日志记录点(journal checkpoint)时的状态
	GetCommittedState(common.Address, common.Hash) common.Hash
	// 返回的是当前VM执行过程中的当前状态(current state),包含未提交的变更
	// 反映了所有最近的写入,但未提交至底层数据库
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	GetTransientState(addr common.Address, key common.Hash) common.Hash
	SetTransientState(addr common.Address, key, val common.Hash)

	SelfDestruct(common.Address)
	HasSelfDestructed(common.Address) bool

	Selfdestruct6780(common.Address)

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for self-destructed accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	AddressInAccessList(addr common.Address) bool
	SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
	// Addcommon.AddressToAccessList adds the given common.Address to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddAddressToAccessList(addr common.Address)
	// AddSlotToAccessList adds the given (common.Address,slot) to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddSlotToAccessList(addr common.Address, slot common.Hash)

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(common.Log)
	AddPreimage(common.Hash, []byte)
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call calls another contract.
	Call(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// CallCode takes another contracts code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// DelegateCall is same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create creates a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error)
}