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

const HashLength int = 32

type Hash [HashLength]byte

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func (h Hash) Bytes() []byte { return h[:] }

var (
	NilHash       Hash = Hash{}
	EmptyCodeHash Hash = Keccak256Hash(nil)
)

func BytesToHash(b []byte) Hash {
	h := Hash{}
	h.SetBytes(b)
	return h
}

func HashToBytes(h Hash) []byte { return h[:] }
