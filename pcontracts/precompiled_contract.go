package pcontracts

import (
	"github.com/lyonnee/evm/common"
	"github.com/lyonnee/evm/params"
)

type PrecompiledContract interface {
	RequiredGas(input []byte) uint64  // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, error) // Run runs the precompiled contract
}

// PrecompiledContractsHomestead contains the default set of pre-compiled Ethereum
// contracts used in the Frontier and Homestead releases.
var PrecompiledContractsHomestead = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}): &ecrecover{},
	common.BytesToAddr([]byte{2}): &sha256hash{},
	common.BytesToAddr([]byte{3}): &ripemd160hash{},
	common.BytesToAddr([]byte{4}): &dataCopy{},
}

// PrecompiledContractsByzantium contains the default set of pre-compiled Ethereum
// contracts used in the Byzantium release.
var PrecompiledContractsByzantium = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}): &ecrecover{},
	common.BytesToAddr([]byte{2}): &sha256hash{},
	common.BytesToAddr([]byte{3}): &ripemd160hash{},
	common.BytesToAddr([]byte{4}): &dataCopy{},
	common.BytesToAddr([]byte{5}): &bigModExp{eip2565: false},
	common.BytesToAddr([]byte{6}): &bn256AddByzantium{},
	common.BytesToAddr([]byte{7}): &bn256ScalarMulByzantium{},
	common.BytesToAddr([]byte{8}): &bn256PairingByzantium{},
}

// PrecompiledContractsIstanbul contains the default set of pre-compiled Ethereum
// contracts used in the Istanbul release.
var PrecompiledContractsIstanbul = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}): &ecrecover{},
	common.BytesToAddr([]byte{2}): &sha256hash{},
	common.BytesToAddr([]byte{3}): &ripemd160hash{},
	common.BytesToAddr([]byte{4}): &dataCopy{},
	common.BytesToAddr([]byte{5}): &bigModExp{eip2565: false},
	common.BytesToAddr([]byte{6}): &bn256AddIstanbul{},
	common.BytesToAddr([]byte{7}): &bn256ScalarMulIstanbul{},
	common.BytesToAddr([]byte{8}): &bn256PairingIstanbul{},
	common.BytesToAddr([]byte{9}): &blake2F{},
}

// PrecompiledContractsBerlin contains the default set of pre-compiled Ethereum
// contracts used in the Berlin release.
var PrecompiledContractsBerlin = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}): &ecrecover{},
	common.BytesToAddr([]byte{2}): &sha256hash{},
	common.BytesToAddr([]byte{3}): &ripemd160hash{},
	common.BytesToAddr([]byte{4}): &dataCopy{},
	common.BytesToAddr([]byte{5}): &bigModExp{eip2565: true},
	common.BytesToAddr([]byte{6}): &bn256AddIstanbul{},
	common.BytesToAddr([]byte{7}): &bn256ScalarMulIstanbul{},
	common.BytesToAddr([]byte{8}): &bn256PairingIstanbul{},
	common.BytesToAddr([]byte{9}): &blake2F{},
}

// PrecompiledContractsCancun contains the default set of pre-compiled Ethereum
// contracts used in the Cancun release.
var PrecompiledContractsCancun = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}):    &ecrecover{},
	common.BytesToAddr([]byte{2}):    &sha256hash{},
	common.BytesToAddr([]byte{3}):    &ripemd160hash{},
	common.BytesToAddr([]byte{4}):    &dataCopy{},
	common.BytesToAddr([]byte{5}):    &bigModExp{eip2565: true},
	common.BytesToAddr([]byte{6}):    &bn256AddIstanbul{},
	common.BytesToAddr([]byte{7}):    &bn256ScalarMulIstanbul{},
	common.BytesToAddr([]byte{8}):    &bn256PairingIstanbul{},
	common.BytesToAddr([]byte{9}):    &blake2F{},
	common.BytesToAddr([]byte{0x0a}): &kzgPointEvaluation{},
}

// PrecompiledContractsBLS contains the set of pre-compiled Ethereum
// contracts specified in EIP-2537. These are exported for testing purposes.
var PrecompiledContractsBLS = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{10}): &bls12381G1Add{},
	common.BytesToAddr([]byte{11}): &bls12381G1Mul{},
	common.BytesToAddr([]byte{12}): &bls12381G1MultiExp{},
	common.BytesToAddr([]byte{13}): &bls12381G2Add{},
	common.BytesToAddr([]byte{14}): &bls12381G2Mul{},
	common.BytesToAddr([]byte{15}): &bls12381G2MultiExp{},
	common.BytesToAddr([]byte{16}): &bls12381Pairing{},
	common.BytesToAddr([]byte{17}): &bls12381MapG1{},
	common.BytesToAddr([]byte{18}): &bls12381MapG2{},
}

var AllPrecompiles = map[common.Address]PrecompiledContract{
	common.BytesToAddr([]byte{1}):    &ecrecover{},
	common.BytesToAddr([]byte{2}):    &sha256hash{},
	common.BytesToAddr([]byte{3}):    &ripemd160hash{},
	common.BytesToAddr([]byte{4}):    &dataCopy{},
	common.BytesToAddr([]byte{5}):    &bigModExp{eip2565: false},
	common.BytesToAddr([]byte{0xf5}): &bigModExp{eip2565: true},
	common.BytesToAddr([]byte{6}):    &bn256AddIstanbul{},
	common.BytesToAddr([]byte{7}):    &bn256ScalarMulIstanbul{},
	common.BytesToAddr([]byte{8}):    &bn256PairingIstanbul{},
	common.BytesToAddr([]byte{9}):    &blake2F{},
	common.BytesToAddr([]byte{0x0a}): &kzgPointEvaluation{},

	common.BytesToAddr([]byte{0x0f, 0x0a}): &bls12381G1Add{},
	common.BytesToAddr([]byte{0x0f, 0x0b}): &bls12381G1Mul{},
	common.BytesToAddr([]byte{0x0f, 0x0c}): &bls12381G1MultiExp{},
	common.BytesToAddr([]byte{0x0f, 0x0d}): &bls12381G2Add{},
	common.BytesToAddr([]byte{0x0f, 0x0e}): &bls12381G2Mul{},
	common.BytesToAddr([]byte{0x0f, 0x0f}): &bls12381G2MultiExp{},
	common.BytesToAddr([]byte{0x0f, 0x10}): &bls12381Pairing{},
	common.BytesToAddr([]byte{0x0f, 0x11}): &bls12381MapG1{},
	common.BytesToAddr([]byte{0x0f, 0x12}): &bls12381MapG2{},
}

var (
	PrecompiledAddressesCancun    []common.Address
	PrecompiledAddressesBerlin    []common.Address
	PrecompiledAddressesIstanbul  []common.Address
	PrecompiledAddressesByzantium []common.Address
	PrecompiledAddressesHomestead []common.Address
)

func init() {
	for k := range PrecompiledContractsHomestead {
		PrecompiledAddressesHomestead = append(PrecompiledAddressesHomestead, k)
	}
	for k := range PrecompiledContractsByzantium {
		PrecompiledAddressesByzantium = append(PrecompiledAddressesByzantium, k)
	}
	for k := range PrecompiledContractsIstanbul {
		PrecompiledAddressesIstanbul = append(PrecompiledAddressesIstanbul, k)
	}
	for k := range PrecompiledContractsBerlin {
		PrecompiledAddressesBerlin = append(PrecompiledAddressesBerlin, k)
	}
	for k := range PrecompiledContractsCancun {
		PrecompiledAddressesCancun = append(PrecompiledAddressesCancun, k)
	}
}

// ActivePrecompiles returns the precompiles enabled with the current configuration.
func ActivePrecompiles(rules params.Rules) []common.Address {
	switch {
	case rules.IsCancun:
		return PrecompiledAddressesCancun
	case rules.IsBerlin:
		return PrecompiledAddressesBerlin
	case rules.IsIstanbul:
		return PrecompiledAddressesIstanbul
	case rules.IsByzantium:
		return PrecompiledAddressesByzantium
	default:
		return PrecompiledAddressesHomestead
	}
}
