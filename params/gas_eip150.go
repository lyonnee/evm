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

package params

const (
	SelfdestructGasEIP150 uint64 = 5000 // Cost of SELFDESTRUCT post EIP 150 (Tangerine)
	CallGasEIP150         uint64 = 700  // Static portion of gas for CALL-derivates after EIP 150 (Tangerine)
	BalanceGasEIP150      uint64 = 400  // The cost of a BALANCE operation after Tangerine
	ExtcodeSizeGasEIP150  uint64 = 700  // Cost of EXTCODESIZE after EIP 150 (Tangerine)
	SloadGasEIP150        uint64 = 200
	ExtcodeCopyBaseEIP150 uint64 = 700
)
