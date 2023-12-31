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

package params

const (
	Keccak256Gas     uint64 = 30 // Once per KECCAK256 operation.
	Keccak256WordGas uint64 = 6  // Once per word of the KECCAK256 operation's data.

	LogGas      uint64 = 375 // Per LOG* operation.
	LogDataGas  uint64 = 8   // Per byte in a LOG* operation's data.
	LogTopicGas uint64 = 375 // Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.

	CreateDataGas uint64 = 200 //

	InitCodeWordGas uint64 = 2 // Once per word of the init code when creating a contract.

	MemoryGas uint64 = 3 // Times the common.Address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL.

	SloadGas         uint64 = 50 // Multiplied by the number of 32-byte words that are copied (round up) for any *COPY operation and added.
	SloadGasFrontier uint64 = 50
	SstoreSetGas     uint64 = 20000 // Once per SSTORE operation.
	SstoreResetGas   uint64 = 5000  // Once per SSTORE operation if the zeroness changes from zero.
	SstoreClearGas   uint64 = 5000  // Once per SSTORE operation if the zeroness doesn't change.
	SstoreRefundGas  uint64 = 15000 // Once per SSTORE operation if the zeroness changes to zero.
	// In EIP-2200: SstoreResetGas was 5000.
	// In EIP-2929: SstoreResetGas was changed to '5000 - COLD_SLOAD_COST'.
	// In EIP-3529: SSTORE_CLEARS_SCHEDULE is defined as SSTORE_RESET_GAS + ACCESS_LIST_STORAGE_KEY_COST
	// Which becomes: 5000 - 2100 + 1900 = 4800
	SstoreClearsScheduleRefundEIP3529 uint64 = SstoreResetGasEIP2200 - ColdSloadCostEIP2929 + TxAccessListStorageKeyGas

	NetSstoreNoopGas  uint64 = 200   // Once per SSTORE operation if the value doesn't change.
	NetSstoreInitGas  uint64 = 20000 // Once per SSTORE operation from clean zero.
	NetSstoreCleanGas uint64 = 5000  // Once per SSTORE operation from clean non-zero.
	NetSstoreDirtyGas uint64 = 200   // Once per SSTORE operation from dirty.

	NetSstoreClearRefund      uint64 = 15000 // Once per SSTORE operation for clearing an originally existing storage slot
	NetSstoreResetRefund      uint64 = 4800  // Once per SSTORE operation for resetting to the original non-zero value
	NetSstoreResetClearRefund uint64 = 19800 // Once per SSTORE operation for resetting to the original zero value

	JumpdestGas uint64 = 1 // Once per JUMPDEST operation.

	// EXP has a dynamic portion depending on the size of the exponent
	ExpByteFrontier              uint64 = 10  // was set to 10 in Frontier
	ExpByteEIP158                uint64 = 50  // was raised to 50 during Eip158 (Spurious Dragon)
	ExpGas                       uint64 = 10  // Once per EXP instruction
	ExpByteGas                   uint64 = 10  // Times ceil(log256(exponent)) for the EXP instruction.
	ExtcodeSizeGasFrontier       uint64 = 20  // Cost of EXTCODESIZE before EIP 150 (Tangerine)
	ExtcodeHashGasConstantinople uint64 = 400 // Cost of EXTCODEHASH (introduced in Constantinople)
	// Extcodecopy has a dynamic AND a static cost. This represents only the
	// static portion of the gas. It was changed during EIP 150 (Tangerine)
	ExtcodeCopyBaseFrontier uint64 = 20

	// CreateBySelfdestructGas is used when the refunded account is one that does
	// not exist. This logic is similar to call.
	// Introduced in Tangerine Whistle (Eip 150)
	CreateBySelfdestructGas uint64 = 25000
	CreateGas               uint64 = 32000 // Once per CREATE operation & contract-creation transaction.
	Create2Gas              uint64 = 32000 // Once per CREATE2 operation

	CopyGas               uint64 = 3     //
	TierStepGas           uint64 = 0     // Once per operation, for a selection of them.
	SelfdestructRefundGas uint64 = 24000 // Refunded following a selfdestruct operation.

	// These have been changed during the course of the chain
	CallGasFrontier      uint64 = 40    // Once per CALL operation & message call transaction.
	CallValueTransferGas uint64 = 9000  // Paid for CALL when the value transfer is non-zero.
	CallNewAccountGas    uint64 = 25000 // Paid for CALL when the destination common.Address didn't exist prior.
	CallStipend          uint64 = 2300  // Free gas given at beginning of call.

	TxGas                     uint64 = 21000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions.
	TxGasContractCreation     uint64 = 53000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions.
	TxDataZeroGas             uint64 = 4     // Per byte of data attached to a transaction that equals zero. NOTE: Not payable on data of calls between transactions.
	TxAccessListStorageKeyGas uint64 = 1900  // Per storage key specified in EIP 2930 access list
	BalanceGasFrontier        uint64 = 20    // The cost of a BALANCE operation

	BlobTxBytesPerFieldElement         = 32      // Size in bytes of a field element
	BlobTxFieldElementsPerBlob         = 4096    // Number of field elements stored in a single data blob
	BlobTxHashVersion                  = 0x01    // Version byte of the commitment hash
	MaxBlobGasPerBlock                 = 1 << 19 // Maximum consumable blob gas for data blobs per block
	BlobTxTargetBlobGasPerBlock        = 1 << 18 // Target consumable blob gas for data blobs per block (for 1559-like pricing)
	BlobTxBlobGasPerBlob               = 1 << 17 // Gas consumption of a single data blob (== blob byte size)
	BlobTxMinBlobGasprice              = 1       // Minimum gas price for data blobs
	BlobTxBlobGaspriceUpdateFraction   = 2225652 // Controls the maximum rate of change for blob gas price
	BlobTxPointEvaluationPrecompileGas = 50000   // Gas price for the point evaluation precompile.
)
