package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/invocarnau/succint-zk-residency/fep-type-1/e2etest/opfaultdisputegame"
	"github.com/urfave/cli/v2"
)

const (
	getGameInfoFlagName = "get-game-info"
	l1RPCFlagName       = "l1-rpc"
	l2RPCFlagName       = "l2-rpc"
	gameFactoryFlagName = "game-factory"
	lastNGamesFlagName  = "last-n-games"
)

func main() {
	app := cli.NewApp()
	app.Name = "OP getter"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    getGameInfoFlagName,
			Aliases: []string{"run", "gi"},
			Usage:   "Get game info for the OP chain",
			Action:  listCommitsToL1,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     l1RPCFlagName,
					Aliases:  []string{"l1"},
					Usage:    "URL of the L1 RPC",
					Required: true,
				},
				&cli.StringFlag{
					Name:     l2RPCFlagName,
					Aliases:  []string{"l2"},
					Usage:    "URL of the L2 RPC",
					Required: true,
				},
				&cli.StringFlag{
					Name:     gameFactoryFlagName,
					Aliases:  []string{"gf", "disputeGameFactoryProxyAddr"},
					Usage:    "Address of the L1 DisputeGameFactoryProxy address",
					Required: true,
				},
				&cli.IntFlag{
					Name:     lastNGamesFlagName,
					Aliases:  []string{"n"},
					Usage:    "Specify the amount of games to list (from last game -n to last game)",
					Required: false,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func listCommitsToL1(cliCtx *cli.Context) error {
	l1URL := cliCtx.String(l1RPCFlagName)
	l2URL := cliCtx.String(l2RPCFlagName)
	disputeGameFactoryProxyAddr := common.HexToAddress(cliCtx.String(gameFactoryFlagName)) // Added for game factory address
	lastNGames := cliCtx.Int(lastNGamesFlagName)                                           // Added for last n games
	fmt.Printf("Getting %d games from %s\n", lastNGames, disputeGameFactoryProxyAddr.Hex())

	clientL1, err := ethclient.Dial(l1URL)
	if err != nil {
		return err
	}
	clientL2, err := ethclient.Dial(l2URL)
	if err != nil {
		return err
	}
	gameFactory, err := bindings.NewDisputeGameFactory(disputeGameFactoryProxyAddr, clientL1)
	if err != nil {
		return err
	}
	gameCounter, err := gameFactory.GameCount(nil)
	if err != nil {
		return err
	}
	for i := int(gameCounter.Int64()) - lastNGames; i < int(gameCounter.Int64()); i++ {
		gameCreation, err := gameFactory.GameAtIndex(nil, big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		gameContract, err := opfaultdisputegame.NewOpfulldisputegame(gameCreation.Proxy, clientL1)
		if err != nil {
			return err
		}

		// check game output
		finalL2BlockNumber, err := gameContract.L2BlockNumber(nil)
		if err != nil {
			return err
		}
		finalRoot, err := gameContract.RootClaim(nil)
		if err != nil {
			return err
		}
		finalOutputV0, err := OutputV0AtBlock(context.Background(), clientL2, finalL2BlockNumber)
		if err != nil {
			return err
		}
		expectedFinalRoot := eth.OutputRoot(finalOutputV0)
		if finalRoot != expectedFinalRoot {
			fmt.Printf("outputV0: %+v\n", finalOutputV0)
			fmt.Printf("final block number: %+v\n", finalL2BlockNumber.Int64())
			fmt.Printf("final claim root: %s\n", common.Hash(finalRoot).Hex())
			fmt.Printf("expectedFinalRoot: %s\n", common.Hash(expectedFinalRoot).Hex())
			return errors.New("RootHash does not match")
		}

		// check game starting root
		initialL2Info, err := gameContract.StartingOutputRoot(nil)
		if err != nil {
			return err
		}
		initialOutputV0, err := OutputV0AtBlock(context.Background(), clientL2, initialL2Info.L2BlockNumber)
		if err != nil {
			return err
		}
		expectedInitialRoot := eth.OutputRoot(initialOutputV0)
		if initialL2Info.Root != expectedInitialRoot {
			fmt.Printf("initial outputV0: %+v\n", initialOutputV0)
			fmt.Printf("initial block number: %+v\n", initialL2Info.L2BlockNumber.Int64())
			fmt.Printf("initial claim root: %s\n", common.Hash(initialL2Info.Root).Hex())
			fmt.Printf("expectedInitialRoot: %s\n", common.Hash(expectedInitialRoot).Hex())
			return errors.New("RootHash does not match")
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
	return nil
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
