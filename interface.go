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
)

type StateDB interface {
	CreateAccount(Address)

	SubBalance(Address, *big.Int)
	AddBalance(Address, *big.Int)
	GetBalance(Address) *big.Int

	GetNonce(Address) uint64
	SetNonce(Address, uint64)

	GetCodeHash(Address) Hash
	GetCode(Address) []byte
	SetCode(Address, []byte)
	GetCodeSize(Address) int

	// 添加退款金额
	AddRefund(uint64)
	// 减去退款金额
	SubRefund(uint64)
	GetRefund() uint64

	// 返回的是当前执行context提交到数据库的已提交状态(committed state)
	// 也就是上次日志记录点(journal checkpoint)时的状态
	GetCommittedState(Address, Hash) Hash
	// 返回的是当前VM执行过程中的当前状态(current state),包含未提交的变更
	// 反映了所有最近的写入,但未提交至底层数据库
	GetState(Address, Hash) Hash
	SetState(Address, Hash, Hash)

	GetTransientState(addr Address, key Hash) Hash
	SetTransientState(addr Address, key, val Hash)

	SelfDestruct(Address)
	HasSelfDestructed(Address) bool

	Selfdestruct6780(Address)

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for self-destructed accounts.
	Exist(Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(Address) bool

	AddressInAccessList(addr Address) bool
	SlotInAccessList(addr Address, slot Hash) (addressOk bool, slotOk bool)
	// Addcommon.AddressToAccessList adds the given define.Address to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddAddressToAccessList(addr Address)
	// AddSlotToAccessList adds the given (define.Address,slot) to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddSlotToAccessList(addr Address, slot Hash)

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(Log)
	AddPreimage(Hash, []byte)
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call calls another contract.
	Call(env *EVM, me ContractRef, addr Address, data []byte, gas, value *big.Int) ([]byte, error)
	// CallCode takes another contracts code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr Address, data []byte, gas, value *big.Int) ([]byte, error)
	// DelegateCall is same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr Address, data []byte, gas *big.Int) ([]byte, error)
	// Create creates a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, Address, error)
}
