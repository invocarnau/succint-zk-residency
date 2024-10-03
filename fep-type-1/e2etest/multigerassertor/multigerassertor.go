// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multigerassertor

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MultigerassertorMetaData contains all meta data concerning the Multigerassertor contract.
var MultigerassertorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"GERs\",\"type\":\"bytes32[]\"}],\"name\":\"CheckGERsExistance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f80fd5b506040516102f93803806102f983398101604081905261002e9161003f565b6001600160a01b031660805261006c565b5f6020828403121561004f575f80fd5b81516001600160a01b0381168114610065575f80fd5b9392505050565b6080516102706100895f395f81816065015260ac01526102705ff3fe608060405234801561000f575f80fd5b5060043610610034575f3560e01c8063cbaf6bbd14610038578063d02103ca14610060575b5f80fd5b61004b61004636600461017c565b61009f565b60405190151581526020015b60405180910390f35b6100877f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610057565b5f805b82811015610170577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663257b36328585848181106100eb576100eb6101eb565b905060200201356040518263ffffffff1660e01b815260040161011091815260200190565b602060405180830381865afa15801561012b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061014f91906101ff565b5f0361015e575f915050610176565b8061016881610216565b9150506100a2565b50600190505b92915050565b5f806020838503121561018d575f80fd5b823567ffffffffffffffff808211156101a4575f80fd5b818501915085601f8301126101b7575f80fd5b8135818111156101c5575f80fd5b8660208260051b85010111156101d9575f80fd5b60209290920196919550909350505050565b634e487b7160e01b5f52603260045260245ffd5b5f6020828403121561020f575f80fd5b5051919050565b5f6001820161023357634e487b7160e01b5f52601160045260245ffd5b506001019056fea2646970667358221220275b83f505e82a86dcb671dcf7758cbf3985ccbf3481c8b09a05d1afad2a893164736f6c63430008140033",
}

// MultigerassertorABI is the input ABI used to generate the binding from.
// Deprecated: Use MultigerassertorMetaData.ABI instead.
var MultigerassertorABI = MultigerassertorMetaData.ABI

// MultigerassertorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MultigerassertorMetaData.Bin instead.
var MultigerassertorBin = MultigerassertorMetaData.Bin

// DeployMultigerassertor deploys a new Ethereum contract, binding an instance of Multigerassertor to it.
func DeployMultigerassertor(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address) (common.Address, *types.Transaction, *Multigerassertor, error) {
	parsed, err := MultigerassertorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultigerassertorBin), backend, _globalExitRootManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Multigerassertor{MultigerassertorCaller: MultigerassertorCaller{contract: contract}, MultigerassertorTransactor: MultigerassertorTransactor{contract: contract}, MultigerassertorFilterer: MultigerassertorFilterer{contract: contract}}, nil
}

// Multigerassertor is an auto generated Go binding around an Ethereum contract.
type Multigerassertor struct {
	MultigerassertorCaller     // Read-only binding to the contract
	MultigerassertorTransactor // Write-only binding to the contract
	MultigerassertorFilterer   // Log filterer for contract events
}

// MultigerassertorCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultigerassertorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultigerassertorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultigerassertorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultigerassertorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultigerassertorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultigerassertorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultigerassertorSession struct {
	Contract     *Multigerassertor // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultigerassertorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultigerassertorCallerSession struct {
	Contract *MultigerassertorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// MultigerassertorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultigerassertorTransactorSession struct {
	Contract     *MultigerassertorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// MultigerassertorRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultigerassertorRaw struct {
	Contract *Multigerassertor // Generic contract binding to access the raw methods on
}

// MultigerassertorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultigerassertorCallerRaw struct {
	Contract *MultigerassertorCaller // Generic read-only contract binding to access the raw methods on
}

// MultigerassertorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultigerassertorTransactorRaw struct {
	Contract *MultigerassertorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultigerassertor creates a new instance of Multigerassertor, bound to a specific deployed contract.
func NewMultigerassertor(address common.Address, backend bind.ContractBackend) (*Multigerassertor, error) {
	contract, err := bindMultigerassertor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multigerassertor{MultigerassertorCaller: MultigerassertorCaller{contract: contract}, MultigerassertorTransactor: MultigerassertorTransactor{contract: contract}, MultigerassertorFilterer: MultigerassertorFilterer{contract: contract}}, nil
}

// NewMultigerassertorCaller creates a new read-only instance of Multigerassertor, bound to a specific deployed contract.
func NewMultigerassertorCaller(address common.Address, caller bind.ContractCaller) (*MultigerassertorCaller, error) {
	contract, err := bindMultigerassertor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultigerassertorCaller{contract: contract}, nil
}

// NewMultigerassertorTransactor creates a new write-only instance of Multigerassertor, bound to a specific deployed contract.
func NewMultigerassertorTransactor(address common.Address, transactor bind.ContractTransactor) (*MultigerassertorTransactor, error) {
	contract, err := bindMultigerassertor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultigerassertorTransactor{contract: contract}, nil
}

// NewMultigerassertorFilterer creates a new log filterer instance of Multigerassertor, bound to a specific deployed contract.
func NewMultigerassertorFilterer(address common.Address, filterer bind.ContractFilterer) (*MultigerassertorFilterer, error) {
	contract, err := bindMultigerassertor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultigerassertorFilterer{contract: contract}, nil
}

// bindMultigerassertor binds a generic wrapper to an already deployed contract.
func bindMultigerassertor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultigerassertorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multigerassertor *MultigerassertorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multigerassertor.Contract.MultigerassertorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multigerassertor *MultigerassertorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multigerassertor.Contract.MultigerassertorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multigerassertor *MultigerassertorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multigerassertor.Contract.MultigerassertorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multigerassertor *MultigerassertorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multigerassertor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multigerassertor *MultigerassertorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multigerassertor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multigerassertor *MultigerassertorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multigerassertor.Contract.contract.Transact(opts, method, params...)
}

// CheckGERsExistance is a free data retrieval call binding the contract method 0xcbaf6bbd.
//
// Solidity: function CheckGERsExistance(bytes32[] GERs) view returns(bool)
func (_Multigerassertor *MultigerassertorCaller) CheckGERsExistance(opts *bind.CallOpts, GERs [][32]byte) (bool, error) {
	var out []interface{}
	err := _Multigerassertor.contract.Call(opts, &out, "CheckGERsExistance", GERs)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckGERsExistance is a free data retrieval call binding the contract method 0xcbaf6bbd.
//
// Solidity: function CheckGERsExistance(bytes32[] GERs) view returns(bool)
func (_Multigerassertor *MultigerassertorSession) CheckGERsExistance(GERs [][32]byte) (bool, error) {
	return _Multigerassertor.Contract.CheckGERsExistance(&_Multigerassertor.CallOpts, GERs)
}

// CheckGERsExistance is a free data retrieval call binding the contract method 0xcbaf6bbd.
//
// Solidity: function CheckGERsExistance(bytes32[] GERs) view returns(bool)
func (_Multigerassertor *MultigerassertorCallerSession) CheckGERsExistance(GERs [][32]byte) (bool, error) {
	return _Multigerassertor.Contract.CheckGERsExistance(&_Multigerassertor.CallOpts, GERs)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Multigerassertor *MultigerassertorCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Multigerassertor.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Multigerassertor *MultigerassertorSession) GlobalExitRootManager() (common.Address, error) {
	return _Multigerassertor.Contract.GlobalExitRootManager(&_Multigerassertor.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Multigerassertor *MultigerassertorCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Multigerassertor.Contract.GlobalExitRootManager(&_Multigerassertor.CallOpts)
}
