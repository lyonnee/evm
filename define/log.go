// Copyright 2023 The evm Authors
// This file is part of the evm library.
//
// The evm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The evm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the evm library. If not, see <http://www.gnu.org/licenses/>.

package define

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
