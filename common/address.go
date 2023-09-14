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
