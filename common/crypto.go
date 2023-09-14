package common

import (
	"hash"

	"github.com/ethereum/go-ethereum/crypto"
)

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// TODO: 可自定义Keccak算法

func NewKeccakState() KeccakState {
	return crypto.NewKeccakState()
}

func Keccak256Hash(data ...[]byte) (h Hash) {
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}
