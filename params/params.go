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

var Bls12381MultiExpDiscountTable = [128]uint64{1200, 888, 764, 641, 594, 547, 500, 453, 438, 423, 408, 394, 379, 364, 349, 334, 330, 326, 322, 318, 314, 310, 306, 302, 298, 294, 289, 285, 281, 277, 273, 269, 268, 266, 265, 263, 262, 260, 259, 257, 256, 254, 253, 251, 250, 248, 247, 245, 244, 242, 241, 239, 238, 236, 235, 233, 232, 231, 229, 228, 226, 225, 223, 222, 221, 220, 219, 219, 218, 217, 216, 216, 215, 214, 213, 213, 212, 211, 211, 210, 209, 208, 208, 207, 206, 205, 205, 204, 203, 202, 202, 201, 200, 199, 199, 198, 197, 196, 196, 195, 194, 193, 193, 192, 191, 191, 190, 189, 188, 188, 187, 186, 185, 185, 184, 183, 182, 182, 181, 180, 179, 179, 178, 177, 176, 176, 175, 174}
