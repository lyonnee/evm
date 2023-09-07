package common

import (
	"github.com/ethereum/go-ethereum/common"
)

// 适配common.Address类型
type Address = common.Address

func BytesToAddr(bs []byte) common.Address {
	return common.Address(bs)
}

func AddrToBytes(addr common.Address) []byte {
	return addr.Bytes()
}

var ZeroAddr common.Address = common.Address{}
