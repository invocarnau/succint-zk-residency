// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package opfaultdisputegame

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

// OpfulldisputegameMetaData contains all meta data concerning the Opfulldisputegame contract.
var OpfulldisputegameMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"GameType\",\"name\":\"_gameType\",\"type\":\"uint32\"},{\"internalType\":\"Claim\",\"name\":\"_absolutePrestate\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_maxGameDepth\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_splitDepth\",\"type\":\"uint256\"},{\"internalType\":\"Duration\",\"name\":\"_gameDuration\",\"type\":\"uint64\"},{\"internalType\":\"contractIBigStepper\",\"name\":\"_vm\",\"type\":\"address\"},{\"internalType\":\"contractIDelayedWETH\",\"name\":\"_weth\",\"type\":\"address\"},{\"internalType\":\"contractIAnchorStateRegistry\",\"name\":\"_anchorStateRegistry\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_l2ChainId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[],\"name\":\"absolutePrestate\",\"outputs\":[{\"internalType\":\"Claim\",\"name\":\"absolutePrestate_\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ident\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_execLeafIdx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_partOffset\",\"type\":\"uint256\"}],\"name\":\"addLocalData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_parentIndex\",\"type\":\"uint256\"},{\"internalType\":\"Claim\",\"name\":\"_claim\",\"type\":\"bytes32\"}],\"name\":\"attack\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"claimCredit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"claimData\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"parentIndex\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"counteredBy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"claimant\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"bond\",\"type\":\"uint128\"},{\"internalType\":\"Claim\",\"name\":\"claim\",\"type\":\"bytes32\"},{\"internalType\":\"Position\",\"name\":\"position\",\"type\":\"uint128\"},{\"internalType\":\"Clock\",\"name\":\"clock\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"claimDataLen\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"len_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createdAt\",\"outputs\":[{\"internalType\":\"Timestamp\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"credit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_parentIndex\",\"type\":\"uint256\"},{\"internalType\":\"Claim\",\"name\":\"_claim\",\"type\":\"bytes32\"}],\"name\":\"defend\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"extraData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"extraData_\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameData\",\"outputs\":[{\"internalType\":\"GameType\",\"name\":\"gameType_\",\"type\":\"uint32\"},{\"internalType\":\"Claim\",\"name\":\"rootClaim_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"extraData_\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameDuration\",\"outputs\":[{\"internalType\":\"Duration\",\"name\":\"gameDuration_\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gameType\",\"outputs\":[{\"internalType\":\"GameType\",\"name\":\"gameType_\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"Position\",\"name\":\"_position\",\"type\":\"uint128\"}],\"name\":\"getRequiredBond\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"requiredBond_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l1Head\",\"outputs\":[{\"internalType\":\"Hash\",\"name\":\"l1Head_\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2BlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"l2BlockNumber_\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"l2ChainId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"l2ChainId_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxGameDepth\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxGameDepth_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_challengeIndex\",\"type\":\"uint256\"},{\"internalType\":\"Claim\",\"name\":\"_claim\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"_isAttack\",\"type\":\"bool\"}],\"name\":\"move\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"enumGameStatus\",\"name\":\"status_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_claimIndex\",\"type\":\"uint256\"}],\"name\":\"resolveClaim\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resolvedAt\",\"outputs\":[{\"internalType\":\"Timestamp\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootClaim\",\"outputs\":[{\"internalType\":\"Claim\",\"name\":\"rootClaim_\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"splitDepth\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"splitDepth_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startingBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"startingBlockNumber_\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startingOutputRoot\",\"outputs\":[{\"internalType\":\"Hash\",\"name\":\"root\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"l2BlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startingRootHash\",\"outputs\":[{\"internalType\":\"Hash\",\"name\":\"startingRootHash_\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"status\",\"outputs\":[{\"internalType\":\"enumGameStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_claimIndex\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_isAttack\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"_stateData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_proof\",\"type\":\"bytes\"}],\"name\":\"step\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vm\",\"outputs\":[{\"internalType\":\"contractIBigStepper\",\"name\":\"vm_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"weth\",\"outputs\":[{\"internalType\":\"contractIDelayedWETH\",\"name\":\"weth_\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"parentIndex\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"Claim\",\"name\":\"claim\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"claimant\",\"type\":\"address\"}],\"name\":\"Move\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"enumGameStatus\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"Resolved\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AlreadyInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AnchorRootNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BondTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotDefendRootClaim\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimAboveSplit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClaimAlreadyResolved\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClockNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClockTimeExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DuplicateStep\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GameDepthExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GameNotInProgress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectBondAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidLocalIdent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPrestate\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidSplitDepth\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoCreditToClaim\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OutOfOrderResolution\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"Claim\",\"name\":\"rootClaim\",\"type\":\"bytes32\"}],\"name\":\"UnexpectedRootClaim\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValidStep\",\"type\":\"error\"}]",
}

// OpfulldisputegameABI is the input ABI used to generate the binding from.
// Deprecated: Use OpfulldisputegameMetaData.ABI instead.
var OpfulldisputegameABI = OpfulldisputegameMetaData.ABI

// Opfulldisputegame is an auto generated Go binding around an Ethereum contract.
type Opfulldisputegame struct {
	OpfulldisputegameCaller     // Read-only binding to the contract
	OpfulldisputegameTransactor // Write-only binding to the contract
	OpfulldisputegameFilterer   // Log filterer for contract events
}

// OpfulldisputegameCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpfulldisputegameCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpfulldisputegameTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpfulldisputegameTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpfulldisputegameFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpfulldisputegameFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpfulldisputegameSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpfulldisputegameSession struct {
	Contract     *Opfulldisputegame // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OpfulldisputegameCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpfulldisputegameCallerSession struct {
	Contract *OpfulldisputegameCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// OpfulldisputegameTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpfulldisputegameTransactorSession struct {
	Contract     *OpfulldisputegameTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// OpfulldisputegameRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpfulldisputegameRaw struct {
	Contract *Opfulldisputegame // Generic contract binding to access the raw methods on
}

// OpfulldisputegameCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpfulldisputegameCallerRaw struct {
	Contract *OpfulldisputegameCaller // Generic read-only contract binding to access the raw methods on
}

// OpfulldisputegameTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpfulldisputegameTransactorRaw struct {
	Contract *OpfulldisputegameTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpfulldisputegame creates a new instance of Opfulldisputegame, bound to a specific deployed contract.
func NewOpfulldisputegame(address common.Address, backend bind.ContractBackend) (*Opfulldisputegame, error) {
	contract, err := bindOpfulldisputegame(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Opfulldisputegame{OpfulldisputegameCaller: OpfulldisputegameCaller{contract: contract}, OpfulldisputegameTransactor: OpfulldisputegameTransactor{contract: contract}, OpfulldisputegameFilterer: OpfulldisputegameFilterer{contract: contract}}, nil
}

// NewOpfulldisputegameCaller creates a new read-only instance of Opfulldisputegame, bound to a specific deployed contract.
func NewOpfulldisputegameCaller(address common.Address, caller bind.ContractCaller) (*OpfulldisputegameCaller, error) {
	contract, err := bindOpfulldisputegame(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpfulldisputegameCaller{contract: contract}, nil
}

// NewOpfulldisputegameTransactor creates a new write-only instance of Opfulldisputegame, bound to a specific deployed contract.
func NewOpfulldisputegameTransactor(address common.Address, transactor bind.ContractTransactor) (*OpfulldisputegameTransactor, error) {
	contract, err := bindOpfulldisputegame(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpfulldisputegameTransactor{contract: contract}, nil
}

// NewOpfulldisputegameFilterer creates a new log filterer instance of Opfulldisputegame, bound to a specific deployed contract.
func NewOpfulldisputegameFilterer(address common.Address, filterer bind.ContractFilterer) (*OpfulldisputegameFilterer, error) {
	contract, err := bindOpfulldisputegame(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpfulldisputegameFilterer{contract: contract}, nil
}

// bindOpfulldisputegame binds a generic wrapper to an already deployed contract.
func bindOpfulldisputegame(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OpfulldisputegameMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Opfulldisputegame *OpfulldisputegameRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Opfulldisputegame.Contract.OpfulldisputegameCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Opfulldisputegame *OpfulldisputegameRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.OpfulldisputegameTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Opfulldisputegame *OpfulldisputegameRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.OpfulldisputegameTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Opfulldisputegame *OpfulldisputegameCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Opfulldisputegame.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Opfulldisputegame *OpfulldisputegameTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Opfulldisputegame *OpfulldisputegameTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.contract.Transact(opts, method, params...)
}

// AbsolutePrestate is a free data retrieval call binding the contract method 0x8d450a95.
//
// Solidity: function absolutePrestate() view returns(bytes32 absolutePrestate_)
func (_Opfulldisputegame *OpfulldisputegameCaller) AbsolutePrestate(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "absolutePrestate")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AbsolutePrestate is a free data retrieval call binding the contract method 0x8d450a95.
//
// Solidity: function absolutePrestate() view returns(bytes32 absolutePrestate_)
func (_Opfulldisputegame *OpfulldisputegameSession) AbsolutePrestate() ([32]byte, error) {
	return _Opfulldisputegame.Contract.AbsolutePrestate(&_Opfulldisputegame.CallOpts)
}

// AbsolutePrestate is a free data retrieval call binding the contract method 0x8d450a95.
//
// Solidity: function absolutePrestate() view returns(bytes32 absolutePrestate_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) AbsolutePrestate() ([32]byte, error) {
	return _Opfulldisputegame.Contract.AbsolutePrestate(&_Opfulldisputegame.CallOpts)
}

// ClaimData is a free data retrieval call binding the contract method 0xc6f0308c.
//
// Solidity: function claimData(uint256 ) view returns(uint32 parentIndex, address counteredBy, address claimant, uint128 bond, bytes32 claim, uint128 position, uint128 clock)
func (_Opfulldisputegame *OpfulldisputegameCaller) ClaimData(opts *bind.CallOpts, arg0 *big.Int) (struct {
	ParentIndex uint32
	CounteredBy common.Address
	Claimant    common.Address
	Bond        *big.Int
	Claim       [32]byte
	Position    *big.Int
	Clock       *big.Int
}, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "claimData", arg0)

	outstruct := new(struct {
		ParentIndex uint32
		CounteredBy common.Address
		Claimant    common.Address
		Bond        *big.Int
		Claim       [32]byte
		Position    *big.Int
		Clock       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ParentIndex = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.CounteredBy = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Claimant = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Bond = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Claim = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.Position = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Clock = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ClaimData is a free data retrieval call binding the contract method 0xc6f0308c.
//
// Solidity: function claimData(uint256 ) view returns(uint32 parentIndex, address counteredBy, address claimant, uint128 bond, bytes32 claim, uint128 position, uint128 clock)
func (_Opfulldisputegame *OpfulldisputegameSession) ClaimData(arg0 *big.Int) (struct {
	ParentIndex uint32
	CounteredBy common.Address
	Claimant    common.Address
	Bond        *big.Int
	Claim       [32]byte
	Position    *big.Int
	Clock       *big.Int
}, error) {
	return _Opfulldisputegame.Contract.ClaimData(&_Opfulldisputegame.CallOpts, arg0)
}

// ClaimData is a free data retrieval call binding the contract method 0xc6f0308c.
//
// Solidity: function claimData(uint256 ) view returns(uint32 parentIndex, address counteredBy, address claimant, uint128 bond, bytes32 claim, uint128 position, uint128 clock)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) ClaimData(arg0 *big.Int) (struct {
	ParentIndex uint32
	CounteredBy common.Address
	Claimant    common.Address
	Bond        *big.Int
	Claim       [32]byte
	Position    *big.Int
	Clock       *big.Int
}, error) {
	return _Opfulldisputegame.Contract.ClaimData(&_Opfulldisputegame.CallOpts, arg0)
}

// ClaimDataLen is a free data retrieval call binding the contract method 0x8980e0cc.
//
// Solidity: function claimDataLen() view returns(uint256 len_)
func (_Opfulldisputegame *OpfulldisputegameCaller) ClaimDataLen(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "claimDataLen")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ClaimDataLen is a free data retrieval call binding the contract method 0x8980e0cc.
//
// Solidity: function claimDataLen() view returns(uint256 len_)
func (_Opfulldisputegame *OpfulldisputegameSession) ClaimDataLen() (*big.Int, error) {
	return _Opfulldisputegame.Contract.ClaimDataLen(&_Opfulldisputegame.CallOpts)
}

// ClaimDataLen is a free data retrieval call binding the contract method 0x8980e0cc.
//
// Solidity: function claimDataLen() view returns(uint256 len_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) ClaimDataLen() (*big.Int, error) {
	return _Opfulldisputegame.Contract.ClaimDataLen(&_Opfulldisputegame.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameCaller) CreatedAt(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "createdAt")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameSession) CreatedAt() (uint64, error) {
	return _Opfulldisputegame.Contract.CreatedAt(&_Opfulldisputegame.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) CreatedAt() (uint64, error) {
	return _Opfulldisputegame.Contract.CreatedAt(&_Opfulldisputegame.CallOpts)
}

// Credit is a free data retrieval call binding the contract method 0xd5d44d80.
//
// Solidity: function credit(address ) view returns(uint256)
func (_Opfulldisputegame *OpfulldisputegameCaller) Credit(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "credit", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Credit is a free data retrieval call binding the contract method 0xd5d44d80.
//
// Solidity: function credit(address ) view returns(uint256)
func (_Opfulldisputegame *OpfulldisputegameSession) Credit(arg0 common.Address) (*big.Int, error) {
	return _Opfulldisputegame.Contract.Credit(&_Opfulldisputegame.CallOpts, arg0)
}

// Credit is a free data retrieval call binding the contract method 0xd5d44d80.
//
// Solidity: function credit(address ) view returns(uint256)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) Credit(arg0 common.Address) (*big.Int, error) {
	return _Opfulldisputegame.Contract.Credit(&_Opfulldisputegame.CallOpts, arg0)
}

// ExtraData is a free data retrieval call binding the contract method 0x609d3334.
//
// Solidity: function extraData() pure returns(bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameCaller) ExtraData(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "extraData")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// ExtraData is a free data retrieval call binding the contract method 0x609d3334.
//
// Solidity: function extraData() pure returns(bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameSession) ExtraData() ([]byte, error) {
	return _Opfulldisputegame.Contract.ExtraData(&_Opfulldisputegame.CallOpts)
}

// ExtraData is a free data retrieval call binding the contract method 0x609d3334.
//
// Solidity: function extraData() pure returns(bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) ExtraData() ([]byte, error) {
	return _Opfulldisputegame.Contract.ExtraData(&_Opfulldisputegame.CallOpts)
}

// GameData is a free data retrieval call binding the contract method 0xfa24f743.
//
// Solidity: function gameData() view returns(uint32 gameType_, bytes32 rootClaim_, bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameCaller) GameData(opts *bind.CallOpts) (struct {
	GameType  uint32
	RootClaim [32]byte
	ExtraData []byte
}, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "gameData")

	outstruct := new(struct {
		GameType  uint32
		RootClaim [32]byte
		ExtraData []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameType = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.RootClaim = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.ExtraData = *abi.ConvertType(out[2], new([]byte)).(*[]byte)

	return *outstruct, err

}

// GameData is a free data retrieval call binding the contract method 0xfa24f743.
//
// Solidity: function gameData() view returns(uint32 gameType_, bytes32 rootClaim_, bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameSession) GameData() (struct {
	GameType  uint32
	RootClaim [32]byte
	ExtraData []byte
}, error) {
	return _Opfulldisputegame.Contract.GameData(&_Opfulldisputegame.CallOpts)
}

// GameData is a free data retrieval call binding the contract method 0xfa24f743.
//
// Solidity: function gameData() view returns(uint32 gameType_, bytes32 rootClaim_, bytes extraData_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) GameData() (struct {
	GameType  uint32
	RootClaim [32]byte
	ExtraData []byte
}, error) {
	return _Opfulldisputegame.Contract.GameData(&_Opfulldisputegame.CallOpts)
}

// GameDuration is a free data retrieval call binding the contract method 0xe1f0c376.
//
// Solidity: function gameDuration() view returns(uint64 gameDuration_)
func (_Opfulldisputegame *OpfulldisputegameCaller) GameDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "gameDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GameDuration is a free data retrieval call binding the contract method 0xe1f0c376.
//
// Solidity: function gameDuration() view returns(uint64 gameDuration_)
func (_Opfulldisputegame *OpfulldisputegameSession) GameDuration() (uint64, error) {
	return _Opfulldisputegame.Contract.GameDuration(&_Opfulldisputegame.CallOpts)
}

// GameDuration is a free data retrieval call binding the contract method 0xe1f0c376.
//
// Solidity: function gameDuration() view returns(uint64 gameDuration_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) GameDuration() (uint64, error) {
	return _Opfulldisputegame.Contract.GameDuration(&_Opfulldisputegame.CallOpts)
}

// GameType is a free data retrieval call binding the contract method 0xbbdc02db.
//
// Solidity: function gameType() view returns(uint32 gameType_)
func (_Opfulldisputegame *OpfulldisputegameCaller) GameType(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "gameType")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GameType is a free data retrieval call binding the contract method 0xbbdc02db.
//
// Solidity: function gameType() view returns(uint32 gameType_)
func (_Opfulldisputegame *OpfulldisputegameSession) GameType() (uint32, error) {
	return _Opfulldisputegame.Contract.GameType(&_Opfulldisputegame.CallOpts)
}

// GameType is a free data retrieval call binding the contract method 0xbbdc02db.
//
// Solidity: function gameType() view returns(uint32 gameType_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) GameType() (uint32, error) {
	return _Opfulldisputegame.Contract.GameType(&_Opfulldisputegame.CallOpts)
}

// GetRequiredBond is a free data retrieval call binding the contract method 0xc395e1ca.
//
// Solidity: function getRequiredBond(uint128 _position) view returns(uint256 requiredBond_)
func (_Opfulldisputegame *OpfulldisputegameCaller) GetRequiredBond(opts *bind.CallOpts, _position *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "getRequiredBond", _position)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRequiredBond is a free data retrieval call binding the contract method 0xc395e1ca.
//
// Solidity: function getRequiredBond(uint128 _position) view returns(uint256 requiredBond_)
func (_Opfulldisputegame *OpfulldisputegameSession) GetRequiredBond(_position *big.Int) (*big.Int, error) {
	return _Opfulldisputegame.Contract.GetRequiredBond(&_Opfulldisputegame.CallOpts, _position)
}

// GetRequiredBond is a free data retrieval call binding the contract method 0xc395e1ca.
//
// Solidity: function getRequiredBond(uint128 _position) view returns(uint256 requiredBond_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) GetRequiredBond(_position *big.Int) (*big.Int, error) {
	return _Opfulldisputegame.Contract.GetRequiredBond(&_Opfulldisputegame.CallOpts, _position)
}

// L1Head is a free data retrieval call binding the contract method 0x6361506d.
//
// Solidity: function l1Head() pure returns(bytes32 l1Head_)
func (_Opfulldisputegame *OpfulldisputegameCaller) L1Head(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "l1Head")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L1Head is a free data retrieval call binding the contract method 0x6361506d.
//
// Solidity: function l1Head() pure returns(bytes32 l1Head_)
func (_Opfulldisputegame *OpfulldisputegameSession) L1Head() ([32]byte, error) {
	return _Opfulldisputegame.Contract.L1Head(&_Opfulldisputegame.CallOpts)
}

// L1Head is a free data retrieval call binding the contract method 0x6361506d.
//
// Solidity: function l1Head() pure returns(bytes32 l1Head_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) L1Head() ([32]byte, error) {
	return _Opfulldisputegame.Contract.L1Head(&_Opfulldisputegame.CallOpts)
}

// L2BlockNumber is a free data retrieval call binding the contract method 0x8b85902b.
//
// Solidity: function l2BlockNumber() pure returns(uint256 l2BlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameCaller) L2BlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "l2BlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L2BlockNumber is a free data retrieval call binding the contract method 0x8b85902b.
//
// Solidity: function l2BlockNumber() pure returns(uint256 l2BlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameSession) L2BlockNumber() (*big.Int, error) {
	return _Opfulldisputegame.Contract.L2BlockNumber(&_Opfulldisputegame.CallOpts)
}

// L2BlockNumber is a free data retrieval call binding the contract method 0x8b85902b.
//
// Solidity: function l2BlockNumber() pure returns(uint256 l2BlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) L2BlockNumber() (*big.Int, error) {
	return _Opfulldisputegame.Contract.L2BlockNumber(&_Opfulldisputegame.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256 l2ChainId_)
func (_Opfulldisputegame *OpfulldisputegameCaller) L2ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "l2ChainId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256 l2ChainId_)
func (_Opfulldisputegame *OpfulldisputegameSession) L2ChainId() (*big.Int, error) {
	return _Opfulldisputegame.Contract.L2ChainId(&_Opfulldisputegame.CallOpts)
}

// L2ChainId is a free data retrieval call binding the contract method 0xd6ae3cd5.
//
// Solidity: function l2ChainId() view returns(uint256 l2ChainId_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) L2ChainId() (*big.Int, error) {
	return _Opfulldisputegame.Contract.L2ChainId(&_Opfulldisputegame.CallOpts)
}

// MaxGameDepth is a free data retrieval call binding the contract method 0xfa315aa9.
//
// Solidity: function maxGameDepth() view returns(uint256 maxGameDepth_)
func (_Opfulldisputegame *OpfulldisputegameCaller) MaxGameDepth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "maxGameDepth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxGameDepth is a free data retrieval call binding the contract method 0xfa315aa9.
//
// Solidity: function maxGameDepth() view returns(uint256 maxGameDepth_)
func (_Opfulldisputegame *OpfulldisputegameSession) MaxGameDepth() (*big.Int, error) {
	return _Opfulldisputegame.Contract.MaxGameDepth(&_Opfulldisputegame.CallOpts)
}

// MaxGameDepth is a free data retrieval call binding the contract method 0xfa315aa9.
//
// Solidity: function maxGameDepth() view returns(uint256 maxGameDepth_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) MaxGameDepth() (*big.Int, error) {
	return _Opfulldisputegame.Contract.MaxGameDepth(&_Opfulldisputegame.CallOpts)
}

// ResolvedAt is a free data retrieval call binding the contract method 0x19effeb4.
//
// Solidity: function resolvedAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameCaller) ResolvedAt(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "resolvedAt")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ResolvedAt is a free data retrieval call binding the contract method 0x19effeb4.
//
// Solidity: function resolvedAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameSession) ResolvedAt() (uint64, error) {
	return _Opfulldisputegame.Contract.ResolvedAt(&_Opfulldisputegame.CallOpts)
}

// ResolvedAt is a free data retrieval call binding the contract method 0x19effeb4.
//
// Solidity: function resolvedAt() view returns(uint64)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) ResolvedAt() (uint64, error) {
	return _Opfulldisputegame.Contract.ResolvedAt(&_Opfulldisputegame.CallOpts)
}

// RootClaim is a free data retrieval call binding the contract method 0xbcef3b55.
//
// Solidity: function rootClaim() pure returns(bytes32 rootClaim_)
func (_Opfulldisputegame *OpfulldisputegameCaller) RootClaim(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "rootClaim")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootClaim is a free data retrieval call binding the contract method 0xbcef3b55.
//
// Solidity: function rootClaim() pure returns(bytes32 rootClaim_)
func (_Opfulldisputegame *OpfulldisputegameSession) RootClaim() ([32]byte, error) {
	return _Opfulldisputegame.Contract.RootClaim(&_Opfulldisputegame.CallOpts)
}

// RootClaim is a free data retrieval call binding the contract method 0xbcef3b55.
//
// Solidity: function rootClaim() pure returns(bytes32 rootClaim_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) RootClaim() ([32]byte, error) {
	return _Opfulldisputegame.Contract.RootClaim(&_Opfulldisputegame.CallOpts)
}

// SplitDepth is a free data retrieval call binding the contract method 0xec5e6308.
//
// Solidity: function splitDepth() view returns(uint256 splitDepth_)
func (_Opfulldisputegame *OpfulldisputegameCaller) SplitDepth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "splitDepth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SplitDepth is a free data retrieval call binding the contract method 0xec5e6308.
//
// Solidity: function splitDepth() view returns(uint256 splitDepth_)
func (_Opfulldisputegame *OpfulldisputegameSession) SplitDepth() (*big.Int, error) {
	return _Opfulldisputegame.Contract.SplitDepth(&_Opfulldisputegame.CallOpts)
}

// SplitDepth is a free data retrieval call binding the contract method 0xec5e6308.
//
// Solidity: function splitDepth() view returns(uint256 splitDepth_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) SplitDepth() (*big.Int, error) {
	return _Opfulldisputegame.Contract.SplitDepth(&_Opfulldisputegame.CallOpts)
}

// StartingBlockNumber is a free data retrieval call binding the contract method 0x70872aa5.
//
// Solidity: function startingBlockNumber() view returns(uint256 startingBlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameCaller) StartingBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "startingBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StartingBlockNumber is a free data retrieval call binding the contract method 0x70872aa5.
//
// Solidity: function startingBlockNumber() view returns(uint256 startingBlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameSession) StartingBlockNumber() (*big.Int, error) {
	return _Opfulldisputegame.Contract.StartingBlockNumber(&_Opfulldisputegame.CallOpts)
}

// StartingBlockNumber is a free data retrieval call binding the contract method 0x70872aa5.
//
// Solidity: function startingBlockNumber() view returns(uint256 startingBlockNumber_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) StartingBlockNumber() (*big.Int, error) {
	return _Opfulldisputegame.Contract.StartingBlockNumber(&_Opfulldisputegame.CallOpts)
}

// StartingOutputRoot is a free data retrieval call binding the contract method 0x57da950e.
//
// Solidity: function startingOutputRoot() view returns(bytes32 root, uint256 l2BlockNumber)
func (_Opfulldisputegame *OpfulldisputegameCaller) StartingOutputRoot(opts *bind.CallOpts) (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "startingOutputRoot")

	outstruct := new(struct {
		Root          [32]byte
		L2BlockNumber *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Root = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.L2BlockNumber = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// StartingOutputRoot is a free data retrieval call binding the contract method 0x57da950e.
//
// Solidity: function startingOutputRoot() view returns(bytes32 root, uint256 l2BlockNumber)
func (_Opfulldisputegame *OpfulldisputegameSession) StartingOutputRoot() (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	return _Opfulldisputegame.Contract.StartingOutputRoot(&_Opfulldisputegame.CallOpts)
}

// StartingOutputRoot is a free data retrieval call binding the contract method 0x57da950e.
//
// Solidity: function startingOutputRoot() view returns(bytes32 root, uint256 l2BlockNumber)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) StartingOutputRoot() (struct {
	Root          [32]byte
	L2BlockNumber *big.Int
}, error) {
	return _Opfulldisputegame.Contract.StartingOutputRoot(&_Opfulldisputegame.CallOpts)
}

// StartingRootHash is a free data retrieval call binding the contract method 0x25fc2ace.
//
// Solidity: function startingRootHash() view returns(bytes32 startingRootHash_)
func (_Opfulldisputegame *OpfulldisputegameCaller) StartingRootHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "startingRootHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StartingRootHash is a free data retrieval call binding the contract method 0x25fc2ace.
//
// Solidity: function startingRootHash() view returns(bytes32 startingRootHash_)
func (_Opfulldisputegame *OpfulldisputegameSession) StartingRootHash() ([32]byte, error) {
	return _Opfulldisputegame.Contract.StartingRootHash(&_Opfulldisputegame.CallOpts)
}

// StartingRootHash is a free data retrieval call binding the contract method 0x25fc2ace.
//
// Solidity: function startingRootHash() view returns(bytes32 startingRootHash_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) StartingRootHash() ([32]byte, error) {
	return _Opfulldisputegame.Contract.StartingRootHash(&_Opfulldisputegame.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_Opfulldisputegame *OpfulldisputegameCaller) Status(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "status")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_Opfulldisputegame *OpfulldisputegameSession) Status() (uint8, error) {
	return _Opfulldisputegame.Contract.Status(&_Opfulldisputegame.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) Status() (uint8, error) {
	return _Opfulldisputegame.Contract.Status(&_Opfulldisputegame.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Opfulldisputegame *OpfulldisputegameCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Opfulldisputegame *OpfulldisputegameSession) Version() (string, error) {
	return _Opfulldisputegame.Contract.Version(&_Opfulldisputegame.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) Version() (string, error) {
	return _Opfulldisputegame.Contract.Version(&_Opfulldisputegame.CallOpts)
}

// Vm is a free data retrieval call binding the contract method 0x3a768463.
//
// Solidity: function vm() view returns(address vm_)
func (_Opfulldisputegame *OpfulldisputegameCaller) Vm(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "vm")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vm is a free data retrieval call binding the contract method 0x3a768463.
//
// Solidity: function vm() view returns(address vm_)
func (_Opfulldisputegame *OpfulldisputegameSession) Vm() (common.Address, error) {
	return _Opfulldisputegame.Contract.Vm(&_Opfulldisputegame.CallOpts)
}

// Vm is a free data retrieval call binding the contract method 0x3a768463.
//
// Solidity: function vm() view returns(address vm_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) Vm() (common.Address, error) {
	return _Opfulldisputegame.Contract.Vm(&_Opfulldisputegame.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address weth_)
func (_Opfulldisputegame *OpfulldisputegameCaller) Weth(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Opfulldisputegame.contract.Call(opts, &out, "weth")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address weth_)
func (_Opfulldisputegame *OpfulldisputegameSession) Weth() (common.Address, error) {
	return _Opfulldisputegame.Contract.Weth(&_Opfulldisputegame.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address weth_)
func (_Opfulldisputegame *OpfulldisputegameCallerSession) Weth() (common.Address, error) {
	return _Opfulldisputegame.Contract.Weth(&_Opfulldisputegame.CallOpts)
}

// AddLocalData is a paid mutator transaction binding the contract method 0xf8f43ff6.
//
// Solidity: function addLocalData(uint256 _ident, uint256 _execLeafIdx, uint256 _partOffset) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) AddLocalData(opts *bind.TransactOpts, _ident *big.Int, _execLeafIdx *big.Int, _partOffset *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "addLocalData", _ident, _execLeafIdx, _partOffset)
}

// AddLocalData is a paid mutator transaction binding the contract method 0xf8f43ff6.
//
// Solidity: function addLocalData(uint256 _ident, uint256 _execLeafIdx, uint256 _partOffset) returns()
func (_Opfulldisputegame *OpfulldisputegameSession) AddLocalData(_ident *big.Int, _execLeafIdx *big.Int, _partOffset *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.AddLocalData(&_Opfulldisputegame.TransactOpts, _ident, _execLeafIdx, _partOffset)
}

// AddLocalData is a paid mutator transaction binding the contract method 0xf8f43ff6.
//
// Solidity: function addLocalData(uint256 _ident, uint256 _execLeafIdx, uint256 _partOffset) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) AddLocalData(_ident *big.Int, _execLeafIdx *big.Int, _partOffset *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.AddLocalData(&_Opfulldisputegame.TransactOpts, _ident, _execLeafIdx, _partOffset)
}

// Attack is a paid mutator transaction binding the contract method 0xc55cd0c7.
//
// Solidity: function attack(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Attack(opts *bind.TransactOpts, _parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "attack", _parentIndex, _claim)
}

// Attack is a paid mutator transaction binding the contract method 0xc55cd0c7.
//
// Solidity: function attack(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Attack(_parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Attack(&_Opfulldisputegame.TransactOpts, _parentIndex, _claim)
}

// Attack is a paid mutator transaction binding the contract method 0xc55cd0c7.
//
// Solidity: function attack(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Attack(_parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Attack(&_Opfulldisputegame.TransactOpts, _parentIndex, _claim)
}

// ClaimCredit is a paid mutator transaction binding the contract method 0x60e27464.
//
// Solidity: function claimCredit(address _recipient) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) ClaimCredit(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "claimCredit", _recipient)
}

// ClaimCredit is a paid mutator transaction binding the contract method 0x60e27464.
//
// Solidity: function claimCredit(address _recipient) returns()
func (_Opfulldisputegame *OpfulldisputegameSession) ClaimCredit(_recipient common.Address) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.ClaimCredit(&_Opfulldisputegame.TransactOpts, _recipient)
}

// ClaimCredit is a paid mutator transaction binding the contract method 0x60e27464.
//
// Solidity: function claimCredit(address _recipient) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) ClaimCredit(_recipient common.Address) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.ClaimCredit(&_Opfulldisputegame.TransactOpts, _recipient)
}

// Defend is a paid mutator transaction binding the contract method 0x35fef567.
//
// Solidity: function defend(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Defend(opts *bind.TransactOpts, _parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "defend", _parentIndex, _claim)
}

// Defend is a paid mutator transaction binding the contract method 0x35fef567.
//
// Solidity: function defend(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Defend(_parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Defend(&_Opfulldisputegame.TransactOpts, _parentIndex, _claim)
}

// Defend is a paid mutator transaction binding the contract method 0x35fef567.
//
// Solidity: function defend(uint256 _parentIndex, bytes32 _claim) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Defend(_parentIndex *big.Int, _claim [32]byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Defend(&_Opfulldisputegame.TransactOpts, _parentIndex, _claim)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Initialize() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Initialize(&_Opfulldisputegame.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Initialize() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Initialize(&_Opfulldisputegame.TransactOpts)
}

// Move is a paid mutator transaction binding the contract method 0x632247ea.
//
// Solidity: function move(uint256 _challengeIndex, bytes32 _claim, bool _isAttack) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Move(opts *bind.TransactOpts, _challengeIndex *big.Int, _claim [32]byte, _isAttack bool) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "move", _challengeIndex, _claim, _isAttack)
}

// Move is a paid mutator transaction binding the contract method 0x632247ea.
//
// Solidity: function move(uint256 _challengeIndex, bytes32 _claim, bool _isAttack) payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Move(_challengeIndex *big.Int, _claim [32]byte, _isAttack bool) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Move(&_Opfulldisputegame.TransactOpts, _challengeIndex, _claim, _isAttack)
}

// Move is a paid mutator transaction binding the contract method 0x632247ea.
//
// Solidity: function move(uint256 _challengeIndex, bytes32 _claim, bool _isAttack) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Move(_challengeIndex *big.Int, _claim [32]byte, _isAttack bool) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Move(&_Opfulldisputegame.TransactOpts, _challengeIndex, _claim, _isAttack)
}

// Resolve is a paid mutator transaction binding the contract method 0x2810e1d6.
//
// Solidity: function resolve() returns(uint8 status_)
func (_Opfulldisputegame *OpfulldisputegameTransactor) Resolve(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "resolve")
}

// Resolve is a paid mutator transaction binding the contract method 0x2810e1d6.
//
// Solidity: function resolve() returns(uint8 status_)
func (_Opfulldisputegame *OpfulldisputegameSession) Resolve() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Resolve(&_Opfulldisputegame.TransactOpts)
}

// Resolve is a paid mutator transaction binding the contract method 0x2810e1d6.
//
// Solidity: function resolve() returns(uint8 status_)
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Resolve() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Resolve(&_Opfulldisputegame.TransactOpts)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0xfdffbb28.
//
// Solidity: function resolveClaim(uint256 _claimIndex) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) ResolveClaim(opts *bind.TransactOpts, _claimIndex *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "resolveClaim", _claimIndex)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0xfdffbb28.
//
// Solidity: function resolveClaim(uint256 _claimIndex) payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) ResolveClaim(_claimIndex *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.ResolveClaim(&_Opfulldisputegame.TransactOpts, _claimIndex)
}

// ResolveClaim is a paid mutator transaction binding the contract method 0xfdffbb28.
//
// Solidity: function resolveClaim(uint256 _claimIndex) payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) ResolveClaim(_claimIndex *big.Int) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.ResolveClaim(&_Opfulldisputegame.TransactOpts, _claimIndex)
}

// Step is a paid mutator transaction binding the contract method 0xd8cc1a3c.
//
// Solidity: function step(uint256 _claimIndex, bool _isAttack, bytes _stateData, bytes _proof) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Step(opts *bind.TransactOpts, _claimIndex *big.Int, _isAttack bool, _stateData []byte, _proof []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.Transact(opts, "step", _claimIndex, _isAttack, _stateData, _proof)
}

// Step is a paid mutator transaction binding the contract method 0xd8cc1a3c.
//
// Solidity: function step(uint256 _claimIndex, bool _isAttack, bytes _stateData, bytes _proof) returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Step(_claimIndex *big.Int, _isAttack bool, _stateData []byte, _proof []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Step(&_Opfulldisputegame.TransactOpts, _claimIndex, _isAttack, _stateData, _proof)
}

// Step is a paid mutator transaction binding the contract method 0xd8cc1a3c.
//
// Solidity: function step(uint256 _claimIndex, bool _isAttack, bytes _stateData, bytes _proof) returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Step(_claimIndex *big.Int, _isAttack bool, _stateData []byte, _proof []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Step(&_Opfulldisputegame.TransactOpts, _claimIndex, _isAttack, _stateData, _proof)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Fallback(&_Opfulldisputegame.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Fallback(&_Opfulldisputegame.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Opfulldisputegame.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Opfulldisputegame *OpfulldisputegameSession) Receive() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Receive(&_Opfulldisputegame.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Opfulldisputegame *OpfulldisputegameTransactorSession) Receive() (*types.Transaction, error) {
	return _Opfulldisputegame.Contract.Receive(&_Opfulldisputegame.TransactOpts)
}

// OpfulldisputegameMoveIterator is returned from FilterMove and is used to iterate over the raw logs and unpacked data for Move events raised by the Opfulldisputegame contract.
type OpfulldisputegameMoveIterator struct {
	Event *OpfulldisputegameMove // Event containing the contract specifics and raw log

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
func (it *OpfulldisputegameMoveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpfulldisputegameMove)
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
		it.Event = new(OpfulldisputegameMove)
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
func (it *OpfulldisputegameMoveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpfulldisputegameMoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpfulldisputegameMove represents a Move event raised by the Opfulldisputegame contract.
type OpfulldisputegameMove struct {
	ParentIndex *big.Int
	Claim       [32]byte
	Claimant    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMove is a free log retrieval operation binding the contract event 0x9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be.
//
// Solidity: event Move(uint256 indexed parentIndex, bytes32 indexed claim, address indexed claimant)
func (_Opfulldisputegame *OpfulldisputegameFilterer) FilterMove(opts *bind.FilterOpts, parentIndex []*big.Int, claim [][32]byte, claimant []common.Address) (*OpfulldisputegameMoveIterator, error) {

	var parentIndexRule []interface{}
	for _, parentIndexItem := range parentIndex {
		parentIndexRule = append(parentIndexRule, parentIndexItem)
	}
	var claimRule []interface{}
	for _, claimItem := range claim {
		claimRule = append(claimRule, claimItem)
	}
	var claimantRule []interface{}
	for _, claimantItem := range claimant {
		claimantRule = append(claimantRule, claimantItem)
	}

	logs, sub, err := _Opfulldisputegame.contract.FilterLogs(opts, "Move", parentIndexRule, claimRule, claimantRule)
	if err != nil {
		return nil, err
	}
	return &OpfulldisputegameMoveIterator{contract: _Opfulldisputegame.contract, event: "Move", logs: logs, sub: sub}, nil
}

// WatchMove is a free log subscription operation binding the contract event 0x9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be.
//
// Solidity: event Move(uint256 indexed parentIndex, bytes32 indexed claim, address indexed claimant)
func (_Opfulldisputegame *OpfulldisputegameFilterer) WatchMove(opts *bind.WatchOpts, sink chan<- *OpfulldisputegameMove, parentIndex []*big.Int, claim [][32]byte, claimant []common.Address) (event.Subscription, error) {

	var parentIndexRule []interface{}
	for _, parentIndexItem := range parentIndex {
		parentIndexRule = append(parentIndexRule, parentIndexItem)
	}
	var claimRule []interface{}
	for _, claimItem := range claim {
		claimRule = append(claimRule, claimItem)
	}
	var claimantRule []interface{}
	for _, claimantItem := range claimant {
		claimantRule = append(claimantRule, claimantItem)
	}

	logs, sub, err := _Opfulldisputegame.contract.WatchLogs(opts, "Move", parentIndexRule, claimRule, claimantRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpfulldisputegameMove)
				if err := _Opfulldisputegame.contract.UnpackLog(event, "Move", log); err != nil {
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

// ParseMove is a log parse operation binding the contract event 0x9b3245740ec3b155098a55be84957a4da13eaf7f14a8bc6f53126c0b9350f2be.
//
// Solidity: event Move(uint256 indexed parentIndex, bytes32 indexed claim, address indexed claimant)
func (_Opfulldisputegame *OpfulldisputegameFilterer) ParseMove(log types.Log) (*OpfulldisputegameMove, error) {
	event := new(OpfulldisputegameMove)
	if err := _Opfulldisputegame.contract.UnpackLog(event, "Move", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OpfulldisputegameResolvedIterator is returned from FilterResolved and is used to iterate over the raw logs and unpacked data for Resolved events raised by the Opfulldisputegame contract.
type OpfulldisputegameResolvedIterator struct {
	Event *OpfulldisputegameResolved // Event containing the contract specifics and raw log

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
func (it *OpfulldisputegameResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OpfulldisputegameResolved)
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
		it.Event = new(OpfulldisputegameResolved)
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
func (it *OpfulldisputegameResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OpfulldisputegameResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OpfulldisputegameResolved represents a Resolved event raised by the Opfulldisputegame contract.
type OpfulldisputegameResolved struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterResolved is a free log retrieval operation binding the contract event 0x5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da60.
//
// Solidity: event Resolved(uint8 indexed status)
func (_Opfulldisputegame *OpfulldisputegameFilterer) FilterResolved(opts *bind.FilterOpts, status []uint8) (*OpfulldisputegameResolvedIterator, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Opfulldisputegame.contract.FilterLogs(opts, "Resolved", statusRule)
	if err != nil {
		return nil, err
	}
	return &OpfulldisputegameResolvedIterator{contract: _Opfulldisputegame.contract, event: "Resolved", logs: logs, sub: sub}, nil
}

// WatchResolved is a free log subscription operation binding the contract event 0x5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da60.
//
// Solidity: event Resolved(uint8 indexed status)
func (_Opfulldisputegame *OpfulldisputegameFilterer) WatchResolved(opts *bind.WatchOpts, sink chan<- *OpfulldisputegameResolved, status []uint8) (event.Subscription, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _Opfulldisputegame.contract.WatchLogs(opts, "Resolved", statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OpfulldisputegameResolved)
				if err := _Opfulldisputegame.contract.UnpackLog(event, "Resolved", log); err != nil {
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

// ParseResolved is a log parse operation binding the contract event 0x5e186f09b9c93491f14e277eea7faa5de6a2d4bda75a79af7a3684fbfb42da60.
//
// Solidity: event Resolved(uint8 indexed status)
func (_Opfulldisputegame *OpfulldisputegameFilterer) ParseResolved(log types.Log) (*OpfulldisputegameResolved, error) {
	event := new(OpfulldisputegameResolved)
	if err := _Opfulldisputegame.contract.UnpackLog(event, "Resolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
