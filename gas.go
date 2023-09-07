package evm

import "github.com/holiman/uint256"

const (
	GasQuickStep   uint64 = 2  // 快速操作的gas价格等级,包括对256位值的算数、位操作、SHA3等
	GasFastestStep uint64 = 3  // 最快操作的gas价格,包括访问stack的操作
	GasFastStep    uint64 = 5  // 快速操作,主要是访问内存的操作
	GasMidStep     uint64 = 8  // 中等速度操作,包括访问存储的操作
	GasSlowStep    uint64 = 10 // 慢速操作,主要包括日志相关操作
	GasExtStep     uint64 = 20 // 更慢的操作,包括外部调用等
)

// callGas 计算调用合约实际消耗的gas成本,包含EIP-150的gas调整规则
func callGas(isEip150 bool, availableGas, base uint64, callCost *uint256.Int) (uint64, error) {
	// 如果启用EIP-150
	if isEip150 {
		// 从可用gas中减去基础gas成本
		availableGas = availableGas - base
		// 再扣除可用gas的63/64作为实际消耗
		gas := availableGas - availableGas/64
		// 如果计算结果超过64位,直接返回可用gas量
		if !callCost.IsUint64() || gas < callCost.Uint64() {
			return gas, nil
		}
	}
	// 如果不启用EIP-150,直接返回调用成本callCost
	if !callCost.IsUint64() {
		// callCost如果超过64位则返回错误
		return 0, ErrGasUintOverflow
	}

	return callCost.Uint64(), nil
}
