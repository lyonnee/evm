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

package define

import (
	"github.com/ethereum/go-ethereum/common"
)

type Hash = common.Hash

// TODO: 修改Hash 类型的定义 和 下面类型转换的 方法

var (
	NilHash       Hash = Hash{}
	EmptyCodeHash Hash = Keccak256Hash(nil)
)

func BytesToHash(b []byte) Hash {
	return common.BytesToHash(b)
}

func HashToBytes(h Hash) []byte {
	return h.Bytes()
}
