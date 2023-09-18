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

package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var NilAddr Address = Address{}

// TODO: 修改Address 类型的定义 和 下面类型转换的 方法

type Address = common.Address

func BytesToAddr(b []byte) Address {
	return common.BytesToAddress(b)
}

func AddrToBytes(a Address) []byte {
	return a.Bytes()
}

// 修改Create方法时,调用处可能也需要修改
var CreateAddress = crypto.CreateAddress
var CreateAddress2 = crypto.CreateAddress2
