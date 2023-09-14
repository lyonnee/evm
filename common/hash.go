package common

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
