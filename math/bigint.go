package math

import "math/big"

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

const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
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
