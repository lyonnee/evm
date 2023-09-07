package evm

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/lyonnee/evm/common"
)

// memoryGasCost计算内存扩展的二次气体
// 只对扩展的内存区域执行此操作，而不是对总内存执行此操作。
func memoryGasCost(mem *Memory, newMemSize uint64) (uint64, error) {
	if newMemSize == 0 {
		return 0, nil
	}

	// uint64中的最大值是max_word_count-1
	// 如果新大小超过可表示的最大值,返回溢出错误
	// 此外，导致newMemSizeWords大于0xFFFFFFFF的newMemSize将导致平方运算溢出。常数0x1FFFFFFFE0是在不溢出气体计算的情况下可以使用的最高数字。
	if newMemSize > 0x1FFFFFFFE0 {
		return 0, ErrGasUintOverflow
	}
	// 计算新内存大小 newMemSize 对应的WordSize newMemSizeWords
	newMemSizeWords := toWordSize(newMemSize)
	newMemSize = newMemSizeWords * 32

	// 如果新大小大于当前内存,则计算扩容的 gas 成本
	if newMemSize > uint64(mem.Len()) {
		square := newMemSizeWords * newMemSizeWords
		linCoef := newMemSizeWords * MemoryGas
		quadCoef := square / QuadCoeffDiv
		newTotalFee := linCoef + quadCoef

		fee := newTotalFee - mem.lastGasCost
		mem.lastGasCost = newTotalFee

		return fee, nil
	}
	return 0, nil
}

/*
	memoryCopierGas,用于计算以下操作码执行时的 gas 消耗:

- CALLDATACOPY
- CODECOPY
- MCOPY
- EXTCODECOPY
- RETURNDATACOPY
*/
func memoryCopierGas(stackpos int) gasFunc {
	return func(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		// 计算扩容内存所需的 gas
		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			return 0, err
		}

		var overflow bool = false
		// 获取要拷贝的数据大小 words
		var wordsSize uint64
		if wordsSize, overflow = stack.Back(stackpos).Uint64WithOverflow(); overflow {
			return 0, ErrGasUintOverflow
		}

		// 对 words 转换为 wordsize 后乘以 CopyGas 得到拷贝数据的 gas
		var copyGasValue uint64
		if copyGasValue, overflow = math.SafeMul(toWordSize(wordsSize), CopyGas); overflow {
			return 0, ErrGasUintOverflow
		}

		// 总 gas 为内存扩容 gas + 拷贝 gas
		var totalGas uint64
		if totalGas, overflow = math.SafeAdd(gas, copyGasValue); overflow {
			return 0, ErrGasUintOverflow
		}
		return totalGas, nil
	}
}

var (
	// CALLDATACOPY (stack position 2)
	gasCallDataCopy gasFunc = memoryCopierGas(2)
	// CODECOPY (stack position 2)
	gasCodeCopy gasFunc = memoryCopierGas(2)
	// MCOPY (stack position 2)
	gasMcopy gasFunc = memoryCopierGas(2)
	// EXTCODECOPY (stack position 3)
	gasExtCodeCopy gasFunc = memoryCopierGas(3)
	// RETURNDATACOPY (stack position 2)
	gasReturnDataCopy gasFunc = memoryCopierGas(2)
)

// gasSStore 实现了以太坊虚拟机(EVM)在 Constantinople 版本中,SSTORE 操作的 gas 计费逻辑
func gasSStore(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	var (
		x, y    = stack.Back(0), stack.Back(1)
		current = evm.StateDB.GetState(contract.Address(), x.Bytes32())
	)

	if evm.chainRules.IsPetersburg || !evm.chainRules.IsConstantinople {
		switch {
		case current == common.ZeroHash && y.Sign() != 0:
			// zero-value -> non-zero value (添加新的值)
			return SstoreSetGas, nil
		case current != common.ZeroHash && y.Sign() == 0:
			// non-zero value -> zero-value (删除值)
			evm.StateDB.AddRefund(SstoreRefundGas)
			return SstoreClearGas, nil
		default:
			// non-zero value -> non-zero value (修改值)
			return SstoreResetGas, nil
		}
	}

	//新的gas计量基于net gas净成本（EIP-1283）

	value := common.Hash(common.BytesToHash(y.Bytes()))
	// 如果当前值等于新值（这是无操作），则扣除200gas。
	if current == value {
		return NetSstoreNoopGas, nil
	}

	// 如果当前值不等于新值

	// 原始值已提交的稳定状态数据
	original := evm.StateDB.GetCommittedState(contract.Address(), common.BytesToHash(x.Bytes()))
	// 如果原始值等于当前值（当前执行上下文未更改此存储槽）
	if original == current {
		if original == common.ZeroHash {
			// 如果原始值为0，则扣除20000gas
			return NetSstoreInitGas, nil
		}
		if value == common.ZeroHash {
			// 若新值为0，则退还 15000gas
			evm.StateDB.AddRefund(NetSstoreClearRefund)
		}
		// 如果原始值不为0，将扣除 5000gas
		return NetSstoreCleanGas, nil
	}

	// 如果原始值不为0
	if original != common.ZeroHash {
		if current == common.ZeroHash {
			// 如果当前值为0（也意味着新值不是0），
			// 则从退款金额中减去 15000gas
			evm.StateDB.SubRefund(NetSstoreClearRefund)
		} else if value == common.ZeroHash {
			// 如果新值为0（也意味着当前值不是0），
			// 则向退款金额添加 15000gas。
			evm.StateDB.SubRefund(NetSstoreClearRefund)
		}
	}

	// 如果原始值等于新值（此存储插槽重置）
	if original == value {
		if original == common.ZeroHash {
			// 如果原始值为0，
			// 则反还 19800gas
			evm.StateDB.AddRefund(NetSstoreResetClearRefund)
		} else {
			// 否则，向退款金额添加 4800gas
			evm.StateDB.AddRefund(NetSstoreResetRefund)
		}
	}
	return NetSstoreDirtyGas, nil
}

func gasSStoreEIP2200(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	if contract.Gas <= SstoreSentryGasEIP2200 {
		// 如果剩余的Gas少于或等于2300，当前调用将失败
		return 0, errors.New("not enough gas for reentrancy sentry")
	}

	var (
		x, y    = stack.Back(0), stack.Back(1)
		current = evm.StateDB.GetState(contract.Address(), x.Bytes32())
	)
	value := common.Hash(common.BytesToHash(y.Bytes()))

	if current == value {
		// 如果当前值等于新值，不执行任何操作，但会扣除SLOAD_GAS。
		return SloadGasEIP2200, nil
	}

	// 如果当前值不等于新值，会根据不同情况扣除Gas，并可能增加或减少退款

	// 原始存储槽的值
	original := evm.StateDB.GetCommittedState(contract.Address(), common.BytesToHash(x.Bytes()))
	// 如果 original 等于 current，表示当前存储槽没有被当前执行上下文修改过
	// 在这种情况下，根据不同的条件扣除Gas，并可能增加或减少退款。
	if original == current {
		if original == common.ZeroHash {
			// 如果 original 为零，表示创建存储槽，扣除 SSTORE_SET_GAS Gas
			return SstoreSetGasEIP2200, nil
		}
		if value == common.ZeroHash {
			// 如果 value 为零，表示删除存储槽，增加 SSTORE_CLEAR_SCHEDULE_REFUND_EIP2200 退款
			evm.StateDB.AddRefund(SstoreClearsScheduleRefundEIP2200)
		}
		return SstoreResetGasEIP2200, nil
	}
	// 如果 original 不等于 current，表示存储槽被修改过（dirty）
	// 在这种情况下，扣除 SLOAD_GAS 并根据不同的条件增加或减少退款

	// 如果 original 不为零
	if original != common.ZeroHash {
		if current == common.ZeroHash {
			// 且 current 为零，表示存储槽从非零值重置为零，扣除 SSTORE_CLEAR_SCHEDULE_REFUND_EIP2200 退款
			evm.StateDB.SubRefund(SstoreClearsScheduleRefundEIP2200)
		} else if value == common.ZeroHash {
			// 如果 value 为零，表示删除存储槽，增加 SSTORE_CLEAR_SCHEDULE_REFUND_EIP2200 退款
			evm.StateDB.AddRefund(SstoreClearsScheduleRefundEIP2200)
		}
	}
	// 如果 original 等于 value，表示存储槽被重置为原始状态
	if original == value {
		if original == common.ZeroHash {
			// 如果 original 为零，表示存储槽从不存在变为存在，增加 SSTORE_SET_GAS - SLOAD_GAS 退款
			evm.StateDB.AddRefund(SstoreSetGasEIP2200 - SloadGasEIP2200)
		} else {
			// 否则，表示存储槽从存在变为存在，增加 SSTORE_RESET_GAS - SLOAD_GAS 退款
			evm.StateDB.AddRefund(SstoreResetGasEIP2200 - SloadGasEIP2200)
		}
	}
	// 返回 SLOAD_GAS Gas，表示存储操作完成
	return SloadGasEIP2200, nil
}

// 计算日志（Log）操作的Gas消耗
// n => 日志数量
func makeGasLog(n uint64) gasFunc {
	return func(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		requestedSize, overflow := stack.Back(1).Uint64WithOverflow()
		if overflow {
			return 0, ErrGasUintOverflow
		}

		gas, err := memoryGasCost(mem, memorySize)
		if err != nil {
			return 0, err
		}

		if gas, overflow = math.SafeAdd(gas, LogGas); overflow {
			return 0, ErrGasUintOverflow
		}
		if gas, overflow = math.SafeAdd(gas, n*LogTopicGas); overflow {
			return 0, ErrGasUintOverflow
		}

		var memorySizeGas uint64
		if memorySizeGas, overflow = math.SafeMul(requestedSize, LogDataGas); overflow {
			return 0, ErrGasUintOverflow
		}
		if gas, overflow = math.SafeAdd(gas, memorySizeGas); overflow {
			return 0, ErrGasUintOverflow
		}

		return gas, nil
	}
}

func gasKeccak256(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	wordGas, overflow := stack.Back(1).Uint64WithOverflow()
	if overflow {
		return 0, ErrGasUintOverflow
	}
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), Keccak256WordGas); overflow {
		return 0, ErrGasUintOverflow
	}
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

// pureMemoryGascost由下面几个操作使用，这些操作除了静态成本外，还有一个动态成本，它完全基于内存扩展
func pureMemoryGascost(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	return memoryGasCost(mem, memorySize)
}

var (
	gasReturn  = pureMemoryGascost
	gasRevert  = pureMemoryGascost
	gasMLoad   = pureMemoryGascost
	gasMStore8 = pureMemoryGascost
	gasMStore  = pureMemoryGascost
	gasCreate  = pureMemoryGascost
)

func gasCreate2(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	wordGas, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow {
		return 0, ErrGasUintOverflow
	}
	if wordGas, overflow = math.SafeMul(toWordSize(wordGas), Keccak256WordGas); overflow {
		return 0, ErrGasUintOverflow
	}
	if gas, overflow = math.SafeAdd(gas, wordGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasCreateEip3860(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	size, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow || size > MaxInitCodeSize {
		return 0, ErrGasUintOverflow
	}
	// Since size <= MaxInitCodeSize, these multiplication cannot overflow
	moreGas := InitCodeWordGas * ((size + 31) / 32)
	if gas, overflow = math.SafeAdd(gas, moreGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasCreate2Eip3860(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	size, overflow := stack.Back(2).Uint64WithOverflow()
	if overflow || size > MaxInitCodeSize {
		return 0, ErrGasUintOverflow
	}
	// Since size <= MaxInitCodeSize, these multiplication cannot overflow
	moreGas := (InitCodeWordGas + Keccak256WordGas) * ((size + 31) / 32)
	if gas, overflow = math.SafeAdd(gas, moreGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasExpFrontier(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	expByteLen := uint64((stack.data[stack.len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * ExpByteFrontier // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = math.SafeAdd(gas, ExpGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasExpEIP158(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	expByteLen := uint64((stack.data[stack.len()-2].BitLen() + 7) / 8)

	var (
		gas      = expByteLen * ExpByteEIP158 // no overflow check required. Max is 256 * ExpByte gas
		overflow bool
	)
	if gas, overflow = math.SafeAdd(gas, ExpGas); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasCall(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	var (
		gas            uint64
		transfersValue = !stack.Back(2).IsZero()
		address        = common.BytesToAddr(stack.Back(1).Bytes())
	)

	if evm.chainRules.IsEIP158 {
		if transfersValue && evm.StateDB.Empty(address) {
			gas += CallNewAccountGas
		}
	} else if !evm.StateDB.Exist(address) {
		gas += CallNewAccountGas
	}
	if transfersValue {
		gas += CallValueTransferGas
	}
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		return 0, ErrGasUintOverflow
	}

	evm.callGasTemp, err = callGas(evm.chainRules.IsEIP150, contract.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = math.SafeAdd(gas, evm.callGasTemp); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasCallCode(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	memoryGas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	var (
		gas      uint64
		overflow bool
	)
	if stack.Back(2).Sign() != 0 {
		gas += CallValueTransferGas
	}
	if gas, overflow = math.SafeAdd(gas, memoryGas); overflow {
		return 0, ErrGasUintOverflow
	}
	evm.callGasTemp, err = callGas(evm.chainRules.IsEIP150, contract.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	if gas, overflow = math.SafeAdd(gas, evm.callGasTemp); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasDelegateCall(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	evm.callGasTemp, err = callGas(evm.chainRules.IsEIP150, contract.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, evm.callGasTemp); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasStaticCall(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	gas, err := memoryGasCost(mem, memorySize)
	if err != nil {
		return 0, err
	}
	evm.callGasTemp, err = callGas(evm.chainRules.IsEIP150, contract.Gas, gas, stack.Back(0))
	if err != nil {
		return 0, err
	}
	var overflow bool
	if gas, overflow = math.SafeAdd(gas, evm.callGasTemp); overflow {
		return 0, ErrGasUintOverflow
	}
	return gas, nil
}

func gasSelfdestruct(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	var gas uint64
	// EIP150 homestead gas reprice fork:
	if evm.chainRules.IsEIP150 {
		gas = SelfdestructGasEIP150
		var address = common.BytesToAddr(stack.Back(0).Bytes())

		if evm.chainRules.IsEIP158 {
			// if empty and transfers value
			if evm.StateDB.Empty(address) && evm.StateDB.GetBalance(contract.Address()).Sign() != 0 {
				gas += CreateBySelfdestructGas
			}
		} else if !evm.StateDB.Exist(address) {
			gas += CreateBySelfdestructGas
		}
	}

	if !evm.StateDB.HasSelfDestructed(contract.Address()) {
		evm.StateDB.AddRefund(SelfdestructRefundGas)
	}
	return gas, nil
}
