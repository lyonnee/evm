package evm

import (
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/lyonnee/evm/common"
	"github.com/lyonnee/evm/math"
	"github.com/lyonnee/evm/pcontracts"
)

type (
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	TransferFunc    func(StateDB, common.Address, common.Address, *big.Int)
	GetHashFunc     func(uint64) common.Hash
)

type BlockContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Block information
	Coinbase      common.Address // Provides information for COINBASE
	GasLimit      uint64         // Provides information for GASLIMIT
	BlockNumber   *big.Int       // Provides information for NUMBER
	Time          uint64         // Provides information for TIME
	Difficulty    *big.Int       // Provides information for DIFFICULTY
	BaseFee       *big.Int       // Provides information for BASEFEE
	Random        *common.Hash   // Provides information for PREVRANDAO
	ExcessBlobGas *uint64        // ExcessBlobGas field in the header, needed to compute the data
}

type TxContext struct {
	// Message information
	Origin     common.Address // Provides information for ORIGIN
	GasPrice   *big.Int       // Provides information for GASPRICE
	BlobHashes []common.Hash  // Provides information for BLOBHASH
}

type codeAndHash struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash) Hash() common.Hash {
	if c.hash == common.ZeroHash {
		c.hash = crypto.Keccak256Hash(c.code)
	}
	return c.hash
}

type EVM struct {
	// Context provides auxiliary blockchain related information
	Context BlockContext
	TxContext
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int
	// chain rules contains the chain rules for the current epoch
	chainRules common.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	Config Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreter *EVMInterpreter
	// abort is used to abort the EVM calling operations
	abort atomic.Bool
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

// Reset resets the EVM with a new transaction context.Reset
// This is not threadsafe and should only be done very cautiously.
func (evm *EVM) Reset(txCtx TxContext, statedb StateDB) {
	evm.TxContext = txCtx
	evm.StateDB = statedb
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	evm.abort.Store(true)
}

// Cancelled returns true if Cancel has been called
func (evm *EVM) Cancelled() bool {
	return evm.abort.Load()
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() *EVMInterpreter {
	return evm.interpreter
}

// SetBlockContext updates the block context of the EVM.
func (evm *EVM) SetBlockContext(blockCtx BlockContext) {
	evm.Context = blockCtx
}

// 调用其他合约
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	// 检查调用深度,避免无限递归调用。
	if evm.depth > int(CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// 检查调用者余额,避免转账不足
	if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	snapshot := evm.StateDB.Snapshot()
	// 根据地址判断是否是预编译合约
	p, isPrecompile := evm.precompile(addr)
	debug := evm.Config.Tracer != nil

	if !evm.StateDB.Exist(addr) {
		// 如果调用的合约不存在,根据EIP158规则判断是否创建合约
		if !isPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
			if debug {
				if evm.depth == 0 {
					evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
					evm.Config.Tracer.CaptureEnd(ret, 0, nil)
				} else {
					evm.Config.Tracer.CaptureEnter(CALL, caller.Address(), addr, input, gas, value)
					evm.Config.Tracer.CaptureExit(ret, 0, nil)
				}
			}
			return nil, gas, nil
		}
		evm.StateDB.CreateAccount(addr)
	}
	// 转账value到被调用合约
	evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)

	// Capture the tracer start/end events in debug mode
	if debug {
		if evm.depth == 0 {
			evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
			defer func(startGas uint64) { // Lazy evaluation of the parameters
				evm.Config.Tracer.CaptureEnd(ret, startGas-gas, err)
			}(gas)
		} else {
			// Handle tracer events for entering and exiting a call frame
			evm.Config.Tracer.CaptureEnter(CALL, caller.Address(), addr, input, gas, value)
			defer func(startGas uint64) {
				evm.Config.Tracer.CaptureExit(ret, startGas-gas, err)
			}(gas)
		}
	}

	if isPrecompile {
		// 如果是预编译合约,直接运行获取返回值
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// 否则,获取代码并创建合约实例,通过解释器Run执行
		code := evm.StateDB.GetCode(addr)
		if len(code) == 0 {
			ret, err = nil, nil
		} else {
			addrCopy := addr
			contract := NewContract(caller, AccountRef(addrCopy), value, gas)
			contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), code)
			ret, err = evm.interpreter.Run(contract, input, false)
			gas = contract.Gas
		}
	}
	if err != nil {
		// 如果执行错误,revert状态到调用前
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.depth > int(CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	snapshot := evm.StateDB.Snapshot()

	// Invoke tracer hooks that signal entering/exiting a call frame
	if evm.Config.Tracer != nil {
		evm.Config.Tracer.CaptureEnter(CALLCODE, caller.Address(), addr, input, gas, value)
		defer func(startGas uint64) {
			evm.Config.Tracer.CaptureExit(ret, startGas-gas, err)
		}(gas)
	}

	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract(caller, AccountRef(caller.Address()), value, gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = evm.interpreter.Run(contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.depth > int(CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	var snapshot = evm.StateDB.Snapshot()

	// Invoke tracer hooks that signal entering/exiting a call frame
	if evm.Config.Tracer != nil {
		// NOTE: caller must, at all times be a contract. It should never happen
		// that caller is something other than a Contract.
		parent := caller.(*Contract)
		// DELEGATECALL inherits value from parent call
		evm.Config.Tracer.CaptureEnter(DELEGATECALL, caller.Address(), addr, input, gas, parent.value)
		defer func(startGas uint64) {
			evm.Config.Tracer.CaptureExit(ret, startGas-gas, err)
		}(gas)
	}

	// It is allowed to call precompiles, even via delegatecall
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// 这里的AsDelegate()为更新了合约的Caller信息
		// caller为上层合约的caller
		contract := NewContract(caller, AccountRef(caller.Address()), nil, gas).AsDelegate()
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = evm.interpreter.Run(contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

// 以给定的输入作为参数执行与addr相关联的合约，
// 但是不允许在调用期间对状态进行任何修改
// 如果试图执行修改状态的操作码将导致异常
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.depth > int(CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	// 可以抛弃的操作,因为修改,所以需要回滚
	// 但这段代码目前没法删除,否则会有异常
	snapshot := evm.StateDB.Snapshot()

	// 调用AddBalance,触发“touch”
	/*	touch
		向账户转账任意数量的以太币,即使是0个。这会导致账户的nonce被重置。
		将数据存储到账户地址所对应的存储空间中。这会重置存储根节点。
		修改账户代码,如通过SELFDESTRUCT清除代码。这会重置代码Hash。
		调用账户的合约,即使静态调用。这会重置账户成为“已初始化”状态。
		通过SUICIDE导致账户被删除。
		在账户首次创建时对其进行初始化。
	*/
	evm.StateDB.AddBalance(addr, math.Big0)

	if evm.Config.Tracer != nil {
		evm.Config.Tracer.CaptureEnter(STATICCALL, caller.Address(), addr, input, gas, nil)
		defer func(startGas uint64) {
			evm.Config.Tracer.CaptureExit(ret, startGas-gas, err)
		}(gas)
	}

	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// At this point, we use a copy of common.Address. If we don't, the go compiler will
		// leak the 'contract' to the outer scope, and make allocation for 'contract'
		// even if the actual execution ends on RunPrecompiled above.
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract(caller, AccountRef(addrCopy), new(big.Int), gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		// When an error was returned by the EVM or when setting the creation code
		// above we revert to the snapshot and consume any gas remaining. Additionally
		// when we're in Homestead this also counts for code storage gas errors.
		ret, err = evm.interpreter.Run(contract, input, true)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = common.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr, CREATE)
}

func (evm *EVM) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *uint256.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = common.CreateAddress2(caller.Address(), salt.Bytes32(), nil)
	return evm.create(caller, &codeAndHash{code: code}, gas, endowment, contractAddr, CREATE2)
}

func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address common.Address, typ OpCode) ([]byte, common.Address, uint64, error) {
	if evm.depth > CALL_CREATE_DEPTH {
		return nil, common.ZeroAddr, gas, ErrDepth
	}
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.ZeroAddr, gas, ErrInsufficientBalance
	}
	nonce := evm.StateDB.GetNonce(caller.Address())
	if nonce+1 < nonce {
		return nil, common.ZeroAddr, gas, ErrNonceUintOverflow
	}
	evm.StateDB.SetNonce(caller.Address(), nonce+1)
	if evm.chainRules.IsBerlin {
		evm.StateDB.AddAddressToAccessList(address)
	}

	contractHash := evm.StateDB.GetCodeHash(address)
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != common.Hash{} || contractHash != EmptyCodeHash) {
		return nil, common.ZeroAddr, 0, ErrContractAddressCollision
	}

	snapshot := evm.StateDB.Snapshot()
	// 创建合约地址的stateObject
	evm.StateDB.CreateAccount(address)
	if evm.chainRules.IsEIP158 {
		evm.StateDB.SetNonce(address, 1)
	}
	// 给合约地址转账,部署合约时可以给合约转账
	evm.Context.Transfer(evm.StateDB, caller.Address(), address, value)

	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.Config.Tracer != nil {
		if evm.depth == 0 {
			evm.Config.Tracer.CaptureStart(evm, caller.Address(), address, true, codeAndHash.code, gas, value)
		} else {
			evm.Config.Tracer.CaptureEnter(typ, caller.Address(), address, codeAndHash.code, gas, value)
		}
	}

	ret, err := evm.interpreter.Run(contract, nil, false)

	if err == nil && evm.chainRules.IsEIP158 && uint64(len(ret)) > MaxCodeSize {
		err = ErrMaxCodeSizeExceeded
	}

	// London - EIP-3541: 拒绝以 0xEF 字节开头的新地址
	if err == nil && len(ret) >= 1 && ret[0] == 0xEF && evm.chainRules.IsLondon {
		err = ErrInvalidCode
	}

	// 如果合约创建成功并且没有返回错误，计算存储代码所需的gas。
	//如果由于没有足够的气体而无法存储代码，则设置一个错误，并让下面的错误检查条件处理它。
	if err == nil {
		createDataGas := uint64(len(ret)) * CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// 当执行上面的创建代码时EVM返回错误，我们将恢复到快照并消耗剩余的gas。
	// ishomestead时，这也会计算代码存储气体错误。
	if err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas) {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	if evm.Config.Tracer != nil {
		if evm.depth == 0 {
			evm.Config.Tracer.CaptureEnd(ret, gas-contract.Gas, err)
		} else {
			evm.Config.Tracer.CaptureExit(ret, gas-contract.Gas, err)
		}
	}
	return ret, address, contract.Gas, err
}

func (evm *EVM) precompile(addr common.Address) (pcontracts.PrecompiledContract, bool) {
	var precompiles map[common.Address]pcontracts.PrecompiledContract
	switch {
	case evm.chainRules.IsCancun:
		precompiles = pcontracts.PrecompiledContractsCancun
	case evm.chainRules.IsBerlin:
		precompiles = pcontracts.PrecompiledContractsBerlin
	case evm.chainRules.IsIstanbul:
		precompiles = pcontracts.PrecompiledContractsIstanbul
	case evm.chainRules.IsByzantium:
		precompiles = pcontracts.PrecompiledContractsByzantium
	default:
		precompiles = pcontracts.PrecompiledContractsHomestead
	}
	p, ok := precompiles[addr]
	return p, ok
}

func NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainId uint64, config Config) *EVM {
	evm := &EVM{
		Context:    blockCtx,
		TxContext:  txCtx,
		StateDB:    statedb,
		Config:     config,
		chainRules: common.NewRules(blockCtx.BlockNumber),
	}
	evm.interpreter = NewEVMInterpreter(evm)
	return evm
}
