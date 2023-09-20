// Copyright 2021 The go-ethereum Authors
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
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lyonnee/evm/define"
	"github.com/lyonnee/evm/params"
)

var loopInterruptTests = []string{
	// infinite loop using JUMP: push(2) jumpdest dup1 jump
	"60025b8056",
	// infinite loop using JUMPI: push(1) push(4) jumpdest dup2 dup2 jumpi
	"600160045b818157",
}

type StateDBImpl struct {
	StateDB
	db *state.StateDB
}

func (s StateDBImpl) CreateAccount(a define.Address) {
	s.db.CreateAccount(a)
}

func (s StateDBImpl) SubBalance(a define.Address, b *big.Int) {
	s.db.SubBalance(a, b)
}
func (s StateDBImpl) AddBalance(a define.Address, b *big.Int) {
	s.db.AddBalance(a, b)
}
func (s StateDBImpl) GetBalance(a define.Address) *big.Int {
	return s.db.GetBalance(a)
}

func (s StateDBImpl) GetNonce(a define.Address) uint64 {
	return s.db.GetNonce(a)
}
func (s StateDBImpl) SetNonce(a define.Address, n uint64) {
	s.db.SetNonce(a, n)
}

func (s StateDBImpl) GetCodeHash(a define.Address) define.Hash {
	return s.db.GetCodeHash(a)
}
func (s StateDBImpl) GetCode(a define.Address) []byte {
	return s.db.GetCode(a)
}
func (s StateDBImpl) SetCode(a define.Address, b []byte) {
	s.db.SetCode(a, b)
}
func (s StateDBImpl) GetCodeSize(a define.Address) int {
	return s.db.GetCodeSize(a)
}

func (s StateDBImpl) AddRefund(n uint64) {
	s.db.AddRefund(n)
}
func (s StateDBImpl) SubRefund(n uint64) {
	s.db.SubRefund(n)
}
func (s StateDBImpl) GetRefund() uint64 {
	return s.db.GetRefund()
}

func (s StateDBImpl) GetCommittedState(a define.Address, k define.Hash) define.Hash {
	return s.db.GetCommittedState(a, k)
}
func (s StateDBImpl) GetState(a define.Address, k define.Hash) define.Hash {
	return s.db.GetState(a, k)
}
func (s StateDBImpl) SetState(a define.Address, k define.Hash, v define.Hash) {
	s.db.SetState(a, k, v)
}

func (s StateDBImpl) GetTransientState(addr define.Address, key define.Hash) define.Hash {
	return s.db.GetTransientState(addr, key)
}
func (s StateDBImpl) SetTransientState(addr define.Address, key, val define.Hash) {
	s.db.SetTransientState(addr, key, val)
}

func (s StateDBImpl) SelfDestruct(a define.Address) {
	s.db.SelfDestruct(a)
}
func (s StateDBImpl) HasSelfDestructed(a define.Address) bool {
	return s.db.HasSelfDestructed(a)
}
func (s StateDBImpl) Selfdestruct6780(a define.Address) {
	s.db.Selfdestruct6780(a)
}

func (s StateDBImpl) Exist(a define.Address) bool {
	return s.db.Exist(a)
}
func (s StateDBImpl) Empty(a define.Address) bool {
	return s.db.Empty(a)
}

func (s StateDBImpl) AddressInAccessList(addr define.Address) bool {
	return s.db.AddressInAccessList(addr)
}
func (s StateDBImpl) SlotInAccessList(addr define.Address, slot define.Hash) (addressOk bool, slotOk bool) {
	return s.db.SlotInAccessList(addr, slot)
}
func (s StateDBImpl) AddAddressToAccessList(addr define.Address) {
	s.db.AddAddressToAccessList(addr)
}
func (s StateDBImpl) AddSlotToAccessList(addr define.Address, slot define.Hash) {
	s.db.AddSlotToAccessList(addr, slot)
}

func (s StateDBImpl) RevertToSnapshot(n int) {
	s.db.RevertToSnapshot(n)
}

func (s StateDBImpl) Snapshot() int {
	return s.db.Snapshot()
}

var allEthashProtocolChanges = &params.ChainConfig{
	ChainID:             big.NewInt(1337),
	HomesteadBlock:      big.NewInt(0),
	DAOForkBlock:        nil,
	DAOForkSupport:      false,
	EIP150Block:         big.NewInt(0),
	EIP155Block:         big.NewInt(0),
	EIP158Block:         big.NewInt(0),
	ByzantiumBlock:      big.NewInt(0),
	ConstantinopleBlock: big.NewInt(0),
	PetersburgBlock:     big.NewInt(0),
	IstanbulBlock:       big.NewInt(0),
	MuirGlacierBlock:    big.NewInt(0),
	BerlinBlock:         big.NewInt(0),
	LondonBlock:         big.NewInt(0),
	ArrowGlacierBlock:   big.NewInt(0),
	GrayGlacierBlock:    big.NewInt(0),
	MergeNetsplitBlock:  nil,
	ShanghaiTime:        nil,
	CancunTime:          nil,
	PragueTime:          nil,
	VerkleTime:          nil,
}

func TestLoopInterrupt(t *testing.T) {
	address := define.BytesToAddr([]byte("contract"))
	vmctx := BlockContext{
		Transfer: func(StateDB, define.Address, define.Address, *big.Int) {},
	}

	for i, tt := range loopInterruptTests {
		statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		statedb.CreateAccount(address)
		statedb.SetCode(address, Hex2Bytes(tt))
		statedb.Finalise(true)

		evm := NewEVM(vmctx, TxContext{}, &StateDBImpl{db: statedb}, allEthashProtocolChanges, Config{})

		errChannel := make(chan error)
		timeout := make(chan bool)

		go func(evm *EVM) {
			_, _, err := evm.Call(AccountRef(define.NilAddr), address, nil, math.MaxUint64, new(big.Int))
			errChannel <- err
		}(evm)

		go func() {
			<-time.After(time.Second)
			timeout <- true
		}()

		evm.Cancel()

		select {
		case <-timeout:
			t.Errorf("test %d timed out", i)
		case err := <-errChannel:
			if err != nil {
				t.Errorf("test %d failure: %v", i, err)
			}
		}
	}
}
