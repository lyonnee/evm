package common

import "github.com/ethereum/go-ethereum/crypto"

type Hash = [32]byte

func BytesToHash(val []byte) Hash {
	return Hash(val)
}

func HashToBytes(hash Hash) []byte {
	return []byte(hash[:])
}

var ZeroHash Hash = Hash{}

var (
	EmptyCodeHash = crypto.Keccak256Hash(nil) // c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
)
