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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lyonnee/evm/params"
)

var loopInterruptTests = []string{
	// infinite loop using JUMP: push(2) jumpdest dup1 jump
	"60025b8056",
	// infinite loop using JUMPI: push(1) push(4) jumpdest dup2 dup2 jumpi
	"600160045b818157",
}

func toGethAddr(a Address) common.Address {
	return common.BytesToAddress(a.Bytes())
}
func toGethHash(h Hash) common.Hash {
	return common.BytesToHash(h.Bytes())
}
func getHashFromGethHash(h common.Hash) Hash {
	return BytesToHash(h.Bytes())
}
func getAddrFromGethAddr(a common.Address) Address {
	return BytesToAddr(a.Bytes())
}

type StateDBImpl struct {
	StateDB
	db *state.StateDB
}

func (s StateDBImpl) CreateAccount(a Address) {
	s.db.CreateAccount(toGethAddr(a))
}

func (s StateDBImpl) SubBalance(a Address, b *big.Int) {
	s.db.SubBalance(toGethAddr(a), b)
}
func (s StateDBImpl) AddBalance(a Address, b *big.Int) {
	s.db.AddBalance(toGethAddr(a), b)
}
func (s StateDBImpl) GetBalance(a Address) *big.Int {
	return s.db.GetBalance(toGethAddr(a))
}

func (s StateDBImpl) GetNonce(a Address) uint64 {
	return s.db.GetNonce(toGethAddr(a))
}
func (s StateDBImpl) SetNonce(a Address, n uint64) {
	s.db.SetNonce(toGethAddr(a), n)
}

func (s StateDBImpl) GetCodeHash(a Address) Hash {
	return getHashFromGethHash(s.db.GetCodeHash(toGethAddr(a)))
}
func (s StateDBImpl) GetCode(a Address) []byte {
	return s.db.GetCode(toGethAddr(a))
}
func (s StateDBImpl) SetCode(a Address, b []byte) {
	s.db.SetCode(toGethAddr(a), b)
}
func (s StateDBImpl) GetCodeSize(a Address) int {
	return s.db.GetCodeSize(toGethAddr(a))
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

func (s StateDBImpl) GetCommittedState(a Address, k Hash) Hash {
	return getHashFromGethHash(s.db.GetCommittedState(toGethAddr(a), common.BytesToHash(k.Bytes())))
}
func (s StateDBImpl) GetState(a Address, k Hash) Hash {
	return getHashFromGethHash(s.db.GetState(toGethAddr(a), common.BytesToHash(k.Bytes())))
}
func (s StateDBImpl) SetState(a Address, k Hash, v Hash) {
	s.db.SetState(toGethAddr(a), common.BytesToHash(k.Bytes()), common.BytesToHash(v.Bytes()))
}

func (s StateDBImpl) GetTransientState(addr Address, key Hash) Hash {
	return getHashFromGethHash(s.db.GetTransientState(toGethAddr(addr), common.BytesToHash(key.Bytes())))
}
func (s StateDBImpl) SetTransientState(addr Address, key, val Hash) {
	s.db.SetTransientState(toGethAddr(addr), common.BytesToHash(key.Bytes()), common.BytesToHash(val.Bytes()))
}

func (s StateDBImpl) SelfDestruct(a Address) {
	s.db.SelfDestruct(toGethAddr(a))
}
func (s StateDBImpl) HasSelfDestructed(a Address) bool {
	return s.db.HasSelfDestructed(toGethAddr(a))
}
func (s StateDBImpl) Selfdestruct6780(a Address) {
	s.db.Selfdestruct6780(toGethAddr(a))
}

func (s StateDBImpl) Exist(a Address) bool {
	return s.db.Exist(toGethAddr(a))
}
func (s StateDBImpl) Empty(a Address) bool {
	return s.db.Empty(toGethAddr(a))
}

func (s StateDBImpl) AddressInAccessList(addr Address) bool {
	return s.db.AddressInAccessList(toGethAddr(addr))
}
func (s StateDBImpl) SlotInAccessList(addr Address, slot Hash) (addressOk bool, slotOk bool) {
	return s.db.SlotInAccessList(toGethAddr(addr), common.BytesToHash(slot.Bytes()))
}
func (s StateDBImpl) AddAddressToAccessList(addr Address) {
	s.db.AddAddressToAccessList(toGethAddr(addr))
}
func (s StateDBImpl) AddSlotToAccessList(addr Address, slot Hash) {
	s.db.AddSlotToAccessList(toGethAddr(addr), common.BytesToHash(slot.Bytes()))
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
	address := BytesToAddr([]byte("contract"))
	vmctx := BlockContext{
		Transfer: func(StateDB, Address, Address, *big.Int) {},
	}

	for i, tt := range loopInterruptTests {
		statedb, _ := state.New(types.EmptyRootHash, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		statedbImpl := &StateDBImpl{db: statedb}
		statedbImpl.CreateAccount(address)
		statedbImpl.SetCode(address, Hex2Bytes(tt))
		statedb.Finalise(true)

		evm := NewEVM(vmctx, TxContext{}, statedbImpl, allEthashProtocolChanges, Config{})

		errChannel := make(chan error)
		timeout := make(chan bool)

		go func(evm *EVM) {
			_, _, err := evm.Call(AccountRef(NilAddr), address, nil, math.MaxUint64, new(big.Int))
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
