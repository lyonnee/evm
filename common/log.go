package common

type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address Address `json:"address"`
	// list of topics provided by the contract.
	Topics []Hash `json:"topics"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash Hash `json:"transactionHash"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex"`
	// hash of the block in which the transaction was included
	BlockHash Hash `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}

func NewLog(addr Address, topics []Hash, data []byte, blockNumber uint64) Log {
	return Log{
		Address:     addr,
		Topics:      topics,
		Data:        data,
		BlockNumber: blockNumber,
	}
}
