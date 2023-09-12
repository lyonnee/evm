package params

const (
	STACK_LIMIT       int = 1024
	CALL_CREATE_DEPTH int = 1024
)

var (
	QuadCoeffDiv    uint64 = 512             // Divisor for the quadratic particle of the memory cost equation.
	MaxCodeSize     uint64 = 24576           // Maximum bytecode to permit for a contract
	MaxInitCodeSize uint64 = 2 * MaxCodeSize // Maximum initcode to permit in a creation transaction and create instructions
	CallCreateDepth uint64 = 1024            // Maximum depth of call/create stack.

)

const (
	StackLimit uint64 = 1024 // Maximum size of VM stack allowed.

	TxAccessListStorageKeyGas uint64 = 1900 // Per storage key specified in EIP 2930 access list

	// In EIP-2200: SstoreResetGas was 5000.
	// In EIP-2929: SstoreResetGas was changed to '5000 - COLD_SLOAD_COST'.
	// In EIP-3529: SSTORE_CLEARS_SCHEDULE is defined as SSTORE_RESET_GAS + ACCESS_LIST_STORAGE_KEY_COST
	// Which becomes: 5000 - 2100 + 1900 = 4800
	SstoreClearsScheduleRefundEIP3529 uint64 = SstoreResetGasEIP2200 - ColdSloadCostEIP2929 + TxAccessListStorageKeyGas

	EpochDuration uint64 = 30000 // Duration between proof-of-work epochs.

)
