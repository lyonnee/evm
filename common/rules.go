package common

import "math/big"

type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsLondon                                      bool
	IsMerge, IsShanghai, IsCancun, IsPrague                 bool
	IsVerkle                                                bool
}

// Rules ensures c's ChainID is not nil.
func NewRules(chainId *big.Int) Rules {
	chainID := chainId
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      true,
		IsEIP150:         true,
		IsEIP155:         true,
		IsEIP158:         true,
		IsByzantium:      true,
		IsConstantinople: true,
		IsPetersburg:     true,
		IsIstanbul:       true,
		IsBerlin:         true,
		IsLondon:         true,
		IsMerge:          true,
		IsShanghai:       true,
		IsCancun:         false,
		IsPrague:         false,
		IsVerkle:         false,
	}
}
