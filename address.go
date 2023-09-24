// Copyright 2023 The evm Authors
// This file is part of the evm library.
//
// The evm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The evm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the evm library. If not, see <http://www.gnu.org/licenses/>.

package evm

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const AddressLength int = 32

type Address [AddressLength]byte

func (a *Address) Bytes() []byte {
	return a[:]
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

var NilAddr Address = Address{}

func BytesToAddr(b []byte) Address {
	a := Address{}
	a.SetBytes(b)
	return a
}

func AddrToBytes(a Address) []byte {
	return a.Bytes()
}

func CreateAddress(a Address, nonce uint64) Address {
	data, _ := rlp.EncodeToBytes([]interface{}{a, nonce})
	return BytesToAddr(crypto.Keccak256(data))
}

func CreateAddress2(a Address, salt [32]byte, inithash []byte) Address {
	return BytesToAddr(crypto.Keccak256([]byte{0xff}, a.Bytes(), salt[:], inithash))
}
