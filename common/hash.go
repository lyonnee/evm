package common

import (
	"github.com/ethereum/go-ethereum/common"
)

type Hash = common.Hash

func BytesToHash(val []byte) common.Hash {
	return common.Hash(val)
}

func HashToBytes(hash common.Hash) []byte {
	return []byte(hash[:])
}

var ZeroHash common.Hash = common.Hash{}
