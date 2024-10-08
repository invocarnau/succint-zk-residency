package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/invocarnau/succint-zk-residency/fep-type-1/e2etest/opfaultdisputegame"
)

const (
	l1URL = "http://localhost:8545"
	l2URL = "http://localhost:9545"
)

var (
	disputeGameFactoryProxyAddr = common.HexToAddress("0x5e8176772863842cd9d692c9E793dc4958626E69")
)

func main() {
	listCommitsToL1()
}

type opHeader struct {
	Hash            common.Hash `json:"hash"`
	StateRoot       common.Hash `json:"stateRoot"`
	WithdrawalsRoot common.Hash `json:"withdrawalsRoot"`
	ParentHash      common.Hash `json:"parentHash"`
}

func listCommitsToL1() {
	clientL1, err := ethclient.Dial(l1URL)
	if err != nil {
		panic(err)
	}
	clientL2, err := ethclient.Dial(l2URL)
	if err != nil {
		panic(err)
	}
	gameFactory, err := bindings.NewDisputeGameFactory(disputeGameFactoryProxyAddr, clientL1)
	if err != nil {
		panic(err)
	}
	gameCounter, err := gameFactory.GameCount(nil)
	if err != nil {
		panic(err)
	}
	for i := 2; i < int(gameCounter.Int64()); i++ {
		gameCreation, err := gameFactory.GameAtIndex(nil, big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}
		gameContract, err := opfaultdisputegame.NewOpfulldisputegame(gameCreation.Proxy, clientL1)
		if err != nil {
			panic(err)
		}

		// check game output
		finalL2BlockNumber, err := gameContract.L2BlockNumber(nil)
		if err != nil {
			panic(err)
		}
		finalRoot, err := gameContract.RootClaim(nil)
		if err != nil {
			panic(err)
		}
		finalOutputV0, err := OutputV0AtBlock(context.Background(), clientL2, finalL2BlockNumber)
		if err != nil {
			panic(err)
		}
		expectedFinalRoot := eth.OutputRoot(finalOutputV0)
		if finalRoot != expectedFinalRoot {
			fmt.Printf("outputV0: %+v\n", finalOutputV0)
			fmt.Printf("final block number: %+v\n", finalL2BlockNumber.Int64())
			fmt.Printf("final claim root: %s\n", common.Hash(finalRoot).Hex())
			fmt.Printf("expectedFinalRoot: %s\n", common.Hash(expectedFinalRoot).Hex())
			panic("RootHash does not match")
		}

		// check game starting root
		initialL2Info, err := gameContract.StartingOutputRoot(nil)
		if err != nil {
			panic(err)
		}
		initialOutputV0, err := OutputV0AtBlock(context.Background(), clientL2, initialL2Info.L2BlockNumber)
		if err != nil {
			panic(err)
		}
		expectedInitialRoot := eth.OutputRoot(initialOutputV0)
		if initialL2Info.Root != expectedInitialRoot {
			fmt.Printf("initial outputV0: %+v\n", initialOutputV0)
			fmt.Printf("initial block number: %+v\n", initialL2Info.L2BlockNumber.Int64())
			fmt.Printf("initial claim root: %s\n", common.Hash(initialL2Info.Root).Hex())
			fmt.Printf("expectedInitialRoot: %s\n", common.Hash(expectedInitialRoot).Hex())
			panic("RootHash does not match")
		}

		fmt.Printf(`
===============================================================================
Game %d
==> Initial
* Block number: %d
* Block hash: %s
* Block root: %s
* MessagePasserStorageRoot: %s
* Claim root: %s
==> Final
* Block number: %d
* Block hash: %s
* Block root: %s
* MessagePasserStorageRoot: %s
* Claim root: %s`,
			i, // game index
			// initial
			initialL2Info.L2BlockNumber.Int64(),                         // Block number
			initialOutputV0.BlockHash.Hex(),                             // Block hash
			common.Hash(initialOutputV0.StateRoot).Hex(),                // Block root
			common.Hash(initialOutputV0.MessagePasserStorageRoot).Hex(), // MessagePasserStorageRoot
			common.Hash(initialL2Info.Root).Hex(),                       // Claim root
			// final
			finalL2BlockNumber.Int64(),                                // Block number
			finalOutputV0.BlockHash.Hex(),                             // Block hash
			common.Hash(finalOutputV0.StateRoot).Hex(),                // Block root
			common.Hash(finalOutputV0.MessagePasserStorageRoot).Hex(), // MessagePasserStorageRoot
			common.Hash(finalRoot).Hex(),                              // Claim root
		)
	}
}

func OutputV0AtBlock(ctx context.Context, client *ethclient.Client, blockNum *big.Int) (*eth.OutputV0, error) {
	header, err := client.HeaderByNumber(ctx, blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get L2 block by hash: %w", err)
	}

	proof, err := GetProof(ctx, client, predeploys.L2ToL1MessagePasserAddr, []common.Hash{}, header.Hash().String())
	if err != nil {
		return nil, fmt.Errorf("failed to get contract proof at block %s: %w", header.Hash(), err)
	}
	if proof == nil {
		return nil, fmt.Errorf("proof %w", ethereum.NotFound)
	}
	// make sure that the proof (including storage hash) that we retrieved is correct by verifying it against the state-root
	if err := proof.Verify(header.Root); err != nil {
		return nil, fmt.Errorf("invalid withdrawal root hash, state root was %s: %w", header.Root, err)
	}
	stateRoot := header.Root
	return &eth.OutputV0{
		StateRoot:                eth.Bytes32(stateRoot),
		MessagePasserStorageRoot: eth.Bytes32(proof.StorageHash),
		BlockHash:                header.Hash(),
	}, nil
}

func GetProof(ctx context.Context, client *ethclient.Client, address common.Address, storage []common.Hash, blockTag string) (*eth.AccountResult, error) {
	var getProofResponse *eth.AccountResult
	err := client.Client().CallContext(ctx, &getProofResponse, "eth_getProof", address, storage, blockTag)
	if err != nil {
		return nil, err
	}
	if getProofResponse == nil {
		return nil, ethereum.NotFound
	}
	if len(getProofResponse.StorageProof) != len(storage) {
		return nil, fmt.Errorf("missing storage proof data, got %d proof entries but requested %d storage keys", len(getProofResponse.StorageProof), len(storage))
	}
	for i, key := range storage {
		if key != getProofResponse.StorageProof[i].Key {
			return nil, fmt.Errorf("unexpected storage proof key difference for entry %d: got %s but requested %s", i, getProofResponse.StorageProof[i].Key, key)
		}
	}
	return getProofResponse, nil
}
