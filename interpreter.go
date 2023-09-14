package evm

import (
	"github.com/lyonnee/evm/common"
	"github.com/lyonnee/evm/math"
)

// Config are the configuration options for the Interpreter
type Config struct {
	Tracer                  EVMLogger // Opcode logger
	NoBaseFee               bool      // Forces the EIP-1559 baseFee to 0 (needed for 0 price calls)
	EnablePreimageRecording bool      // Enables recording of SHA3/keccak preimages
	ExtraEips               []int     // Additional EIPS that are to be enabled
}

// ScopeContext contains the things that are per-call, such as stack and memory,
// but not transients like pc and gas
type ScopeContext struct {
	Memory   *Memory
	Stack    *Stack
	Contract *Contract
}

// EVMInterpreter represents an EVM interpreter
type EVMInterpreter struct {
	evm   *EVM
	table *JumpTable

	hasher    common.KeccakState // Keccak256 hasher instance shared across opcodes
	hasherBuf common.Hash        // Keccak256 hasher result array shared aross opcodes

	readOnly   bool   // Whether to throw on stateful modifications
	returnData []byte // Last CALL's return data for subsequent reuse
}

// NewEVMInterpreter returns a new instance of the Interpreter.
func NewEVMInterpreter(evm *EVM) *EVMInterpreter {
	// If jump table was not initialised we set the default one.
	var table *JumpTable
	switch {
	case evm.chainRules.IsCancun:
		table = &cancunInstructionSet
	case evm.chainRules.IsShanghai:
		table = &shanghaiInstructionSet
	case evm.chainRules.IsMerge:
		table = &mergeInstructionSet
	case evm.chainRules.IsLondon:
		table = &londonInstructionSet
	case evm.chainRules.IsBerlin:
		table = &berlinInstructionSet
	case evm.chainRules.IsIstanbul:
		table = &istanbulInstructionSet
	case evm.chainRules.IsConstantinople:
		table = &constantinopleInstructionSet
	case evm.chainRules.IsByzantium:
		table = &byzantiumInstructionSet
	case evm.chainRules.IsEIP158:
		table = &spuriousDragonInstructionSet
	case evm.chainRules.IsEIP150:
		table = &tangerineWhistleInstructionSet
	case evm.chainRules.IsHomestead:
		table = &homesteadInstructionSet
	default:
		table = &frontierInstructionSet
	}
	var extraEips []int
	if len(evm.Config.ExtraEips) > 0 {
		// Deep-copy jumptable to prevent modification of opcodes in other tables
		table = copyJumpTable(table)
	}
	for _, eip := range evm.Config.ExtraEips {
		if err := EnableEIP(eip, table); err != nil {
			// Disable it, so caller can check if it's activated or not
			// log.Error("EIP activation failed", "eip", eip, "error", err)
		} else {
			extraEips = append(extraEips, eip)
		}
	}
	evm.Config.ExtraEips = extraEips
	return &EVMInterpreter{evm: evm, table: table}
}

func (in *EVMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	// 调用深度+1,限制最大为1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}
	// 重置上个调用的返回数据,因为每个返回的调用都会有新的返回数据
	in.returnData = nil
	// 如果没有代码,不执行
	if len(contract.Code) == 0 {
		return nil, nil
	}

	var (
		op          OpCode           // 当前操作码
		mem         = NewMemory()    // 分配内存
		stack       = newstack()     // 局部栈
		callContext = &ScopeContext{ // 调用上下文
			Memory:   mem,
			Stack:    stack,
			Contract: contract,
		}
		// For optimisation reason we're using uint64 as the program counter.
		// It's theoretically possible to go above 2^64. The YP defines the PC
		// to be uint256. Practically much less so feasible.
		// 为优化使用uint64作为程序计数器,理论上可能超过2^64,黄皮书中定义为uint256
		pc   = uint64(0) // 程序计数器
		cost uint64
		// copies used by tracer
		pcCopy  uint64                        // 追踪需要的拷贝
		gasCopy uint64                        // 追踪执行前的gas剩余
		logged  bool                          // 是否已经记录日志
		res     []byte                        // 操作结果
		debug   = in.evm.Config.Tracer != nil // 是否为DEBUG模式
	)

	defer func() {
		returnStack(stack)
	}()
	contract.Input = input

	if debug {
		defer func() {
			if err != nil {
				if !logged {
					in.evm.Config.Tracer.CaptureState(pcCopy, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
				} else {
					in.evm.Config.Tracer.CaptureFault(pcCopy, op, gasCopy, cost, callContext, in.evm.depth, err)
				}
			}
		}()
	}
	// 主循环,直到触发STOP/RETURN/SELFDESTRUCT或者错误
	for {
		if debug {
			logged, pcCopy, gasCopy = false, pc, contract.Gas
		}

		op = contract.GetOp(pc)
		operation := in.table[op]
		cost = operation.constantGas

		if sLen := stack.len(); sLen < operation.minStack {
			return nil, &ErrStackUnderflow{
				stackLen: sLen,
				required: operation.minStack,
			}
		} else if sLen > operation.maxStack {
			return nil, &ErrStackOverflow{
				stackLen: sLen,
				limit:    operation.maxStack,
			}
		}
		if !contract.UseGas(cost) {
			return nil, ErrOutOfGas
		}
		if operation.dynamicGas != nil {
			var memorySize uint64

			if operation.memorySize != nil {
				memSize, overflow := operation.memorySize(stack)
				if overflow {
					return nil, ErrGasUintOverflow
				}
				if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
					return nil, ErrGasUintOverflow
				}
			}

			var dynamicCost uint64
			dynamicCost, err = operation.dynamicGas(in.evm, contract, stack, mem, memorySize)
			cost += dynamicCost
			if err != nil || !contract.UseGas(dynamicCost) {
				return nil, ErrOutOfGas
			}
			if debug {
				in.evm.Config.Tracer.CaptureState(pc, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
				logged = true
			}
			if memorySize > 0 {
				mem.Resize(memorySize)
			}
		} else if debug {
			in.evm.Config.Tracer.CaptureState(pc, op, gasCopy, cost, callContext, in.returnData, in.evm.depth, err)
			logged = true
		}

		res, err = operation.execute(&pc, in, callContext)
		if err != nil {
			break
		}
		pc++
	}

	if err == errStopToken {
		err = nil // clear stop token error
	}

	return res, err
}
