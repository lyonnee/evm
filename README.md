# EVM

A reusable Ethereum Virtual Machine implementation in Go.

fork[go-ethereum](https://github.com/ethereum/go-ethereum) Mawinor (v1.12.2).

## Introduction

EVM is a Go package that provides a modular implementation of the Ethereum Virtual Machine, by abstracting dependencies and making it easy to integrate into any Ethereum-compatible blockchain.

It is based on the core EVM code from [go-ethereum](https://github.com/ethereum/go-ethereum), but designed to be decoupled from other Ethereum specifics like `State` and `BlockChain` interfaces.

### Features

- Modular design, minimal dependencies
- Compatible with Ethereum EVM specifications
- Implements opcodes, gas costs, stack, memory, etc.
- Easy integration with custom `State` implementations 

## Contributing

Contributions are welcome! Open an issue or PR.

## License

EVM is released under the [GNU GENERAL PUBLIC LICENSE](LICENSE).
