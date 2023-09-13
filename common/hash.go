package common

import "github.com/ethereum/go-ethereum/crypto"

const HashLength int = 32

type Hash [HashLength]byte

func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

func HashToBytes(hash Hash) []byte {
	return []byte(hash[:])
}

var ZeroHash Hash = Hash{}

var (
	EmptyCodeHash = [32]byte(crypto.Keccak256Hash(nil)) // c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
)
