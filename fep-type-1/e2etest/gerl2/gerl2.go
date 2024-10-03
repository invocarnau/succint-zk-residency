// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gerl2

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

// Gerl2MetaData contains all meta data concerning the Gerl2 contract.
var Gerl2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"GlobalExitRootAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAllowedContracts\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyCoinbase\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newGlobalExitRoot\",\"type\":\"bytes32\"}],\"name\":\"InsertGlobalExitRoot\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_newRoot\",\"type\":\"bytes32\"}],\"name\":\"insertGlobalExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561000f575f80fd5b5060405161031b38038061031b83398101604081905261002e9161003f565b6001600160a01b031660805261006c565b5f6020828403121561004f575f80fd5b81516001600160a01b0381168114610065575f80fd5b9392505050565b60805161029161008a5f395f818160d001526101e801526102915ff3fe608060405234801561000f575f80fd5b5060043610610064575f3560e01c8063257b36321161004d578063257b36321461009957806333d6247d146100b8578063a3c573eb146100cb575f80fd5b806301fd90441461006857806312da06b214610084575b5f80fd5b61007160015481565b6040519081526020015b60405180910390f35b610097610092366004610244565b610117565b005b6100716100a7366004610244565b5f6020819052908152604090205481565b6100976100c6366004610244565b6101d0565b6100f27f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161007b565b413314610150576040517f116c64a800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f81815260208190526040812054900361019e575f818152602081905260408082204290555182917fb1b866fe5fac68e8f1a4ab2520c7a6b493a954934bbd0f054bd91d6674a4c0d591a250565b6040517f1f97a58200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000161461023f576040517fb49365dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600155565b5f60208284031215610254575f80fd5b503591905056fea2646970667358221220732d86ccfdb7a4e919b79043bf2cb087d41ebb78a6838191a10d16ce253c4f6564736f6c63430008140033",
}

// Gerl2ABI is the input ABI used to generate the binding from.
// Deprecated: Use Gerl2MetaData.ABI instead.
var Gerl2ABI = Gerl2MetaData.ABI

// Gerl2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Gerl2MetaData.Bin instead.
var Gerl2Bin = Gerl2MetaData.Bin

// DeployGerl2 deploys a new Ethereum contract, binding an instance of Gerl2 to it.
func DeployGerl2(auth *bind.TransactOpts, backend bind.ContractBackend, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Gerl2, error) {
	parsed, err := Gerl2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Gerl2Bin), backend, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Gerl2{Gerl2Caller: Gerl2Caller{contract: contract}, Gerl2Transactor: Gerl2Transactor{contract: contract}, Gerl2Filterer: Gerl2Filterer{contract: contract}}, nil
}

// Gerl2 is an auto generated Go binding around an Ethereum contract.
type Gerl2 struct {
	Gerl2Caller     // Read-only binding to the contract
	Gerl2Transactor // Write-only binding to the contract
	Gerl2Filterer   // Log filterer for contract events
}

// Gerl2Caller is an auto generated read-only Go binding around an Ethereum contract.
type Gerl2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Gerl2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Gerl2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Gerl2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Gerl2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Gerl2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Gerl2Session struct {
	Contract     *Gerl2            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Gerl2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Gerl2CallerSession struct {
	Contract *Gerl2Caller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Gerl2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Gerl2TransactorSession struct {
	Contract     *Gerl2Transactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Gerl2Raw is an auto generated low-level Go binding around an Ethereum contract.
type Gerl2Raw struct {
	Contract *Gerl2 // Generic contract binding to access the raw methods on
}

// Gerl2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Gerl2CallerRaw struct {
	Contract *Gerl2Caller // Generic read-only contract binding to access the raw methods on
}

// Gerl2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Gerl2TransactorRaw struct {
	Contract *Gerl2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewGerl2 creates a new instance of Gerl2, bound to a specific deployed contract.
func NewGerl2(address common.Address, backend bind.ContractBackend) (*Gerl2, error) {
	contract, err := bindGerl2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Gerl2{Gerl2Caller: Gerl2Caller{contract: contract}, Gerl2Transactor: Gerl2Transactor{contract: contract}, Gerl2Filterer: Gerl2Filterer{contract: contract}}, nil
}

// NewGerl2Caller creates a new read-only instance of Gerl2, bound to a specific deployed contract.
func NewGerl2Caller(address common.Address, caller bind.ContractCaller) (*Gerl2Caller, error) {
	contract, err := bindGerl2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Gerl2Caller{contract: contract}, nil
}

// NewGerl2Transactor creates a new write-only instance of Gerl2, bound to a specific deployed contract.
func NewGerl2Transactor(address common.Address, transactor bind.ContractTransactor) (*Gerl2Transactor, error) {
	contract, err := bindGerl2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Gerl2Transactor{contract: contract}, nil
}

// NewGerl2Filterer creates a new log filterer instance of Gerl2, bound to a specific deployed contract.
func NewGerl2Filterer(address common.Address, filterer bind.ContractFilterer) (*Gerl2Filterer, error) {
	contract, err := bindGerl2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Gerl2Filterer{contract: contract}, nil
}

// bindGerl2 binds a generic wrapper to an already deployed contract.
func bindGerl2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Gerl2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gerl2 *Gerl2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gerl2.Contract.Gerl2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gerl2 *Gerl2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gerl2.Contract.Gerl2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gerl2 *Gerl2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gerl2.Contract.Gerl2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gerl2 *Gerl2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gerl2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gerl2 *Gerl2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gerl2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gerl2 *Gerl2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gerl2.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Gerl2 *Gerl2Caller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Gerl2.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Gerl2 *Gerl2Session) BridgeAddress() (common.Address, error) {
	return _Gerl2.Contract.BridgeAddress(&_Gerl2.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Gerl2 *Gerl2CallerSession) BridgeAddress() (common.Address, error) {
	return _Gerl2.Contract.BridgeAddress(&_Gerl2.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Gerl2 *Gerl2Caller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Gerl2.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Gerl2 *Gerl2Session) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Gerl2.Contract.GlobalExitRootMap(&_Gerl2.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Gerl2 *Gerl2CallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Gerl2.Contract.GlobalExitRootMap(&_Gerl2.CallOpts, arg0)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Gerl2 *Gerl2Caller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Gerl2.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Gerl2 *Gerl2Session) LastRollupExitRoot() ([32]byte, error) {
	return _Gerl2.Contract.LastRollupExitRoot(&_Gerl2.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Gerl2 *Gerl2CallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Gerl2.Contract.LastRollupExitRoot(&_Gerl2.CallOpts)
}

// InsertGlobalExitRoot is a paid mutator transaction binding the contract method 0x12da06b2.
//
// Solidity: function insertGlobalExitRoot(bytes32 _newRoot) returns()
func (_Gerl2 *Gerl2Transactor) InsertGlobalExitRoot(opts *bind.TransactOpts, _newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.contract.Transact(opts, "insertGlobalExitRoot", _newRoot)
}

// InsertGlobalExitRoot is a paid mutator transaction binding the contract method 0x12da06b2.
//
// Solidity: function insertGlobalExitRoot(bytes32 _newRoot) returns()
func (_Gerl2 *Gerl2Session) InsertGlobalExitRoot(_newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.Contract.InsertGlobalExitRoot(&_Gerl2.TransactOpts, _newRoot)
}

// InsertGlobalExitRoot is a paid mutator transaction binding the contract method 0x12da06b2.
//
// Solidity: function insertGlobalExitRoot(bytes32 _newRoot) returns()
func (_Gerl2 *Gerl2TransactorSession) InsertGlobalExitRoot(_newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.Contract.InsertGlobalExitRoot(&_Gerl2.TransactOpts, _newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Gerl2 *Gerl2Transactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Gerl2 *Gerl2Session) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.Contract.UpdateExitRoot(&_Gerl2.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Gerl2 *Gerl2TransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Gerl2.Contract.UpdateExitRoot(&_Gerl2.TransactOpts, newRoot)
}

// Gerl2InsertGlobalExitRootIterator is returned from FilterInsertGlobalExitRoot and is used to iterate over the raw logs and unpacked data for InsertGlobalExitRoot events raised by the Gerl2 contract.
type Gerl2InsertGlobalExitRootIterator struct {
	Event *Gerl2InsertGlobalExitRoot // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Gerl2InsertGlobalExitRootIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Gerl2InsertGlobalExitRoot)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Gerl2InsertGlobalExitRoot)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Gerl2InsertGlobalExitRootIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Gerl2InsertGlobalExitRootIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Gerl2InsertGlobalExitRoot represents a InsertGlobalExitRoot event raised by the Gerl2 contract.
type Gerl2InsertGlobalExitRoot struct {
	NewGlobalExitRoot [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterInsertGlobalExitRoot is a free log retrieval operation binding the contract event 0xb1b866fe5fac68e8f1a4ab2520c7a6b493a954934bbd0f054bd91d6674a4c0d5.
//
// Solidity: event InsertGlobalExitRoot(bytes32 indexed newGlobalExitRoot)
func (_Gerl2 *Gerl2Filterer) FilterInsertGlobalExitRoot(opts *bind.FilterOpts, newGlobalExitRoot [][32]byte) (*Gerl2InsertGlobalExitRootIterator, error) {

	var newGlobalExitRootRule []interface{}
	for _, newGlobalExitRootItem := range newGlobalExitRoot {
		newGlobalExitRootRule = append(newGlobalExitRootRule, newGlobalExitRootItem)
	}

	logs, sub, err := _Gerl2.contract.FilterLogs(opts, "InsertGlobalExitRoot", newGlobalExitRootRule)
	if err != nil {
		return nil, err
	}
	return &Gerl2InsertGlobalExitRootIterator{contract: _Gerl2.contract, event: "InsertGlobalExitRoot", logs: logs, sub: sub}, nil
}

// WatchInsertGlobalExitRoot is a free log subscription operation binding the contract event 0xb1b866fe5fac68e8f1a4ab2520c7a6b493a954934bbd0f054bd91d6674a4c0d5.
//
// Solidity: event InsertGlobalExitRoot(bytes32 indexed newGlobalExitRoot)
func (_Gerl2 *Gerl2Filterer) WatchInsertGlobalExitRoot(opts *bind.WatchOpts, sink chan<- *Gerl2InsertGlobalExitRoot, newGlobalExitRoot [][32]byte) (event.Subscription, error) {

	var newGlobalExitRootRule []interface{}
	for _, newGlobalExitRootItem := range newGlobalExitRoot {
		newGlobalExitRootRule = append(newGlobalExitRootRule, newGlobalExitRootItem)
	}

	logs, sub, err := _Gerl2.contract.WatchLogs(opts, "InsertGlobalExitRoot", newGlobalExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Gerl2InsertGlobalExitRoot)
				if err := _Gerl2.contract.UnpackLog(event, "InsertGlobalExitRoot", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInsertGlobalExitRoot is a log parse operation binding the contract event 0xb1b866fe5fac68e8f1a4ab2520c7a6b493a954934bbd0f054bd91d6674a4c0d5.
//
// Solidity: event InsertGlobalExitRoot(bytes32 indexed newGlobalExitRoot)
func (_Gerl2 *Gerl2Filterer) ParseInsertGlobalExitRoot(log types.Log) (*Gerl2InsertGlobalExitRoot, error) {
	event := new(Gerl2InsertGlobalExitRoot)
	if err := _Gerl2.contract.UnpackLog(event, "InsertGlobalExitRoot", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
