package params

const (
	SstoreSentryGasEIP2200            uint64 = 2300  // Minimum gas required to be present for an SSTORE call, not consumed
	SstoreSetGasEIP2200               uint64 = 20000 // Once per SSTORE operation from clean zero to non-zero
	SstoreResetGasEIP2200             uint64 = 5000  // Once per SSTORE operation from clean non-zero to something else
	SstoreClearsScheduleRefundEIP2200 uint64 = 15000 // Once per SSTORE operation for clearing an originally existing storage slot
	SloadGasEIP2200                   uint64 = 800   // Cost of SLOAD after EIP 2200 (part of Istanbul)
)
