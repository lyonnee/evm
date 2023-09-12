package common

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// 适配common.Address类型
const AddressLength int = 32

type Address [AddressLength]byte

func (a Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

func BytesToAddr(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func AddrToBytes(addr Address) []byte {
	return addr[:]
}

var ZeroAddr Address = Address{}

func CreateAddress(b Address, nonce uint64) Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return BytesToAddr(crypto.Keccak256(data)[12:])
}

func CreateAddress2(b Address, salt [32]byte, inithash []byte) Address {
	return BytesToAddr(crypto.Keccak256([]byte{0xff}, b[:], salt[:], inithash)[12:])
}
