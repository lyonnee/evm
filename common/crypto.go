package common

import (
	"github.com/ethereum/go-ethereum/crypto"
)

var CreateAddress func(addr Address, nonce uint64) Address = crypto.CreateAddress

var CreateAddress2 func(addr Address, salt [32]byte, inithash []byte) Address = crypto.CreateAddress2
