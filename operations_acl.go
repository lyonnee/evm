package evm

import (
	"errors"

	"github.com/lyonnee/evm/common"
	"github.com/lyonnee/evm/math"
	"github.com/lyonnee/evm/params"
)

func makeGasSStoreFunc(clearingRefund uint64) gasFunc {
	return func(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		if contract.Gas <= params.SstoreSentryGasEIP2200 {
			return 0, errors.New("not enough gas for reentrancy sentry")
		}

		var (
			x, y    = stack.peek(), stack.Back(1)
			slot    = common.BytesToHash(x.Bytes())
			current = evm.StateDB.GetState(contract.Address(), slot)
			cost    = uint64(0)
		)

		if addrPresent, slotPresent := evm.StateDB.SlotInAccessList(contract.Address(), slot); !slotPresent {
			cost = params.ColdSloadCostEIP2929
			// If the caller cannot afford the cost, this change will be rolled back
			evm.StateDB.AddSlotToAccessList(contract.Address(), slot)
			if !addrPresent {
				// Once we're done with YOLOv2 and schedule this for mainnet, might
				// be good to remove this panic here, which is just really a
				// canary to have during testing
				panic("impossible case: common.Address was not present in access list during sstore op")
			}
		}
		value := common.BytesToHash(y.Bytes())

		if current == value {
			return cost + params.WarmStorageReadCostEIP2929, nil
		}
		original := evm.StateDB.GetCommittedState(contract.Address(), common.BytesToHash(x.Bytes()))
		if original == current {
			if original == common.NilHash {
				return cost + params.SstoreSetGasEIP2200, nil
			}
			if value == common.NilHash {
				evm.StateDB.AddRefund(clearingRefund)
			}
			return cost + (params.SstoreResetGasEIP2200 - params.ColdSloadCostEIP2929), nil
		}
		if original != common.NilHash {
			if current == common.NilHash {
				evm.StateDB.SubRefund(clearingRefund)
			} else if value == common.NilHash {
				evm.StateDB.AddRefund(clearingRefund)
			}
		}
		if original == value {
			if original == common.NilHash {
				evm.StateDB.AddRefund(params.SstoreSetGasEIP2200 - params.WarmStorageReadCostEIP2929)
			} else {
				evm.StateDB.AddRefund((params.SstoreResetGasEIP2200 - params.ColdSloadCostEIP2929) - params.WarmStorageReadCostEIP2929)
			}
		}
		return cost + params.WarmStorageReadCostEIP2929, nil
	}
}

func gasSLoadEIP2929(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	loc := stack.peek()
	slot := common.BytesToHash(loc.Bytes())

	if _, slotPresent := evm.StateDB.SlotInAccessList(contract.Address(), slot); !slotPresent {
		evm.StateDB.AddSlotToAccessList(contract.Address(), slot)
		return params.ColdSloadCostEIP2929, nil
	}
	return params.WarmStorageReadCostEIP2929, nil
}

func gasExtCodeCopyEIP2929(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	// memory expansion first (dynamic part of pre-2929 implementation)
	gas, err := gasExtCodeCopy(evm, contract, stack, mem, memorySize)
	if err != nil {
		return 0, err
	}
	addr := common.BytesToAddr(stack.peek().Bytes())
	// Check slot presence in the access list
	if !evm.StateDB.AddressInAccessList(addr) {
		evm.StateDB.AddAddressToAccessList(addr)
		var overflow bool
		// We charge (cold-warm), since 'warm' is already charged as constantGas
		if gas, overflow = math.SafeAdd(gas, params.ColdAccountAccessCostEIP2929-params.WarmStorageReadCostEIP2929); overflow {
			return 0, ErrGasUintOverflow
		}
		return gas, nil
	}
	return gas, nil
}

func gasEip2929AccountCheck(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
	addr := common.BytesToAddr(stack.peek().Bytes())
	// Check slot presence in the access list
	if !evm.StateDB.AddressInAccessList(addr) {
		// If the caller cannot afford the cost, this change will be rolled back
		evm.StateDB.AddAddressToAccessList(addr)
		// The warm storage read cost is already charged as constantGas
		return params.ColdAccountAccessCostEIP2929 - params.WarmStorageReadCostEIP2929, nil
	}
	return 0, nil
}

func makeCallVariantGasCallEIP2929(oldCalculator gasFunc) gasFunc {
	return func(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		addr := common.BytesToAddr(stack.peek().Bytes())
		// Check slot presence in the access list
		warmAccess := evm.StateDB.AddressInAccessList(addr)
		// The WarmStorageReadCostEIP2929 (100) is already deducted in the form of a constant cost, so
		// the cost to charge for cold access, if any, is Cold - Warm
		coldCost := params.ColdAccountAccessCostEIP2929 - params.WarmStorageReadCostEIP2929
		if !warmAccess {
			evm.StateDB.AddAddressToAccessList(addr)
			// Charge the remaining difference here already, to correctly calculate available
			// gas for call
			if !contract.UseGas(coldCost) {
				return 0, ErrOutOfGas
			}
		}
		// Now call the old calculator, which takes into account
		// - create new account
		// - transfer value
		// - memory expansion
		// - 63/64ths rule
		gas, err := oldCalculator(evm, contract, stack, mem, memorySize)
		if warmAccess || err != nil {
			return gas, err
		}
		// In case of a cold access, we temporarily add the cold charge back, and also
		// add it to the returned gas. By adding it to the return, it will be charged
		// outside of this function, as part of the dynamic gas, and that will make it
		// also become correctly reported to tracers.
		contract.Gas += coldCost
		return gas + coldCost, nil
	}
}

var (
	gasCallEIP2929         = makeCallVariantGasCallEIP2929(gasCall)
	gasDelegateCallEIP2929 = makeCallVariantGasCallEIP2929(gasDelegateCall)
	gasStaticCallEIP2929   = makeCallVariantGasCallEIP2929(gasStaticCall)
	gasCallCodeEIP2929     = makeCallVariantGasCallEIP2929(gasCallCode)
	gasSelfdestructEIP2929 = makeSelfdestructGasFn(true)
	// gasSelfdestructEIP3529 implements the changes in EIP-2539 (no refunds)
	gasSelfdestructEIP3529 = makeSelfdestructGasFn(false)

	// gasSStoreEIP2929 implements gas cost for SSTORE according to EIP-2929
	//
	// When calling SSTORE, check if the (common.Address, storage_key) pair is in accessed_storage_keys.
	// If it is not, charge an additional COLD_SLOAD_COST gas, and add the pair to accessed_storage_keys.
	// Additionally, modify the parameters defined in EIP 2200 as follows:
	//
	// Parameter 	Old value 	New value
	// SLOAD_GAS 	800 	= WARM_STORAGE_READ_COST
	// SSTORE_RESET_GAS 	5000 	5000 - COLD_SLOAD_COST
	//
	//The other parameters defined in EIP 2200 are unchanged.
	// see gasSStoreEIP2200(...) in core/vm/gas_table.go for more info about how EIP 2200 is specified
	gasSStoreEIP2929 = makeGasSStoreFunc(params.SstoreClearsScheduleRefundEIP2200)

	// gasSStoreEIP2539 implements gas cost for SSTORE according to EIP-2539
	// Replace `SSTORE_CLEARS_SCHEDULE` with `SSTORE_RESET_GAS + ACCESS_LIST_STORAGE_KEY_COST` (4,800)
	gasSStoreEIP3529 = makeGasSStoreFunc(params.SstoreClearsScheduleRefundEIP3529)
)

func makeSelfdestructGasFn(refundsEnabled bool) gasFunc {
	gasFunc := func(evm *EVM, contract *Contract, stack *Stack, mem *Memory, memorySize uint64) (uint64, error) {
		var (
			gas     uint64
			address = common.BytesToAddr(stack.peek().Bytes())
		)
		if !evm.StateDB.AddressInAccessList(address) {
			// If the caller cannot afford the cost, this change will be rolled back
			evm.StateDB.AddAddressToAccessList(address)
			gas = params.ColdAccountAccessCostEIP2929
		}
		// if empty and transfers value
		if evm.StateDB.Empty(address) && evm.StateDB.GetBalance(contract.Address()).Sign() != 0 {
			gas += params.CreateBySelfdestructGas
		}
		if refundsEnabled && !evm.StateDB.HasSelfDestructed(contract.Address()) {
			evm.StateDB.AddRefund(params.SelfdestructRefundGas)
		}
		return gas, nil
	}
	return gasFunc
}
