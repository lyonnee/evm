package params

const (
	SelfdestructGasEIP150 uint64 = 5000 // Cost of SELFDESTRUCT post EIP 150 (Tangerine)
	CallGasEIP150         uint64 = 700  // Static portion of gas for CALL-derivates after EIP 150 (Tangerine)
	BalanceGasEIP150      uint64 = 400  // The cost of a BALANCE operation after Tangerine
	ExtcodeSizeGasEIP150  uint64 = 700  // Cost of EXTCODESIZE after EIP 150 (Tangerine)
	SloadGasEIP150        uint64 = 200
	ExtcodeCopyBaseEIP150 uint64 = 700
)
