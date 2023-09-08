package math

import (
	"math/big"
)

var (
	Big0      = big.NewInt(0)
	Big1      = big.NewInt(1)
	Big3      = big.NewInt(3)
	Big4      = big.NewInt(4)
	Big7      = big.NewInt(7)
	Big8      = big.NewInt(8)
	Big16     = big.NewInt(16)
	Big20     = big.NewInt(20)
	Big32     = big.NewInt(32)
	Big64     = big.NewInt(64)
	Big96     = big.NewInt(96)
	Big480    = big.NewInt(480)
	Big1024   = big.NewInt(1024)
	Big3072   = big.NewInt(3072)
	Big199680 = big.NewInt(199680)
)

func BigMax(first, second *big.Int) *big.Int {
	if BigGreaterThan(first, second) {
		return new(big.Int).Set(first)
	} else {
		return new(big.Int).Set(second)
	}
}

func BigGreaterThan(first, second *big.Int) bool {
	return first.Cmp(second) > 0
}
