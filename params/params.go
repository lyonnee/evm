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

	EpochDuration uint64 = 30000 // Duration between proof-of-work epochs.
)
