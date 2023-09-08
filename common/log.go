package common

import "github.com/ethereum/go-ethereum/core/types"

type Loger interface {
	UnmarshalJSON([]byte) error
	MarshalJSON() ([]byte, error)
}

type Log = types.Log

func NewLog(addr Address, topics []Hash, data []byte, blockNumber uint64) Loger {
	return &Log{
		Address:     addr,
		Topics:      topics,
		Data:        data,
		BlockNumber: blockNumber,
	}
}
