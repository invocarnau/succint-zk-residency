package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/invocarnau/succint-zk-residency/bridge"
	"github.com/invocarnau/succint-zk-residency/goutils"
	"github.com/invocarnau/succint-zk-residency/goutils/contracts/opfaultdisputegame"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	proof                  = "proof"
	l1ChainIDFlagName      = "l1-chainid"
	l2ChainIDFlagName      = "l2-chainid"
	gameFactoryFlagName    = "game-factory"
	gerL1FlagName          = "ger-l1"
	gerL2FlagName          = "ger-l2"
	prevL2BlockNumFlagName = "prev-l2-block-num"
	networkIDFlagName      = "network-id"
)

func main() {
	app := cli.NewApp()
	app.Name = "OP getter"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    proof,
			Aliases: []string{},
			Usage:   "Generate a chain proof for an OP stack chain",
			Action:  generateChainProof,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     l1ChainIDFlagName,
					Aliases:  []string{"l1", "l1-chain-id"},
					Usage:    "Chain ID of the L1 network",
					Required: true,
				},
				&cli.StringFlag{
					Name:     l2ChainIDFlagName,
					Aliases:  []string{"l2", "l2-chain-id"},
					Usage:    "Chain ID of the L2 network",
					Required: true,
				},
				&cli.StringFlag{
					Name:     gameFactoryFlagName,
					Aliases:  []string{"gf", "disputeGameFactoryProxyAddr"},
					Usage:    "Address of the L1 DisputeGameFactoryProxy contract",
					Required: true,
				},
				&cli.IntFlag{
					Name:     prevL2BlockNumFlagName,
					Aliases:  []string{"from", "prev-block", "pb"},
					Usage:    "Specify the previous L2 block number, from which the proof will start, this should match the new block of the last proof",
					Required: true,
				},
				&cli.StringFlag{
					Name:     gerL1FlagName,
					Aliases:  []string{"gl1", "ger-addr-l1"},
					Usage:    "Address of the L1 GER contract",
					Required: true,
				},
				&cli.StringFlag{
					Name:     gerL2FlagName,
					Aliases:  []string{"gl2", "ger-addr-l2"},
					Usage:    "Address of the L2 GER contract",
					Required: true,
				},
				&cli.IntFlag{
					Name:     networkIDFlagName,
					Aliases:  []string{"nid", "network"},
					Usage:    "Specify the network ID of the rollup (according to the rollup manager)",
					Required: true,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func generateChainProof(cliCtx *cli.Context) error {
	op, err := NewOP(cliCtx)
	if err != nil {
		return fmt.Errorf("failed to call NewOP: %w", err)
	}
	return op.GenerateChainProof(cliCtx.Context)
}

type OP struct {
	l1ChainID       int
	l2ChainID       int
	networkID       uint64
	l1RPC           *ethclient.Client
	l2RPC           *ethclient.Client
	gameFactoryAddr common.Address
	gameFactory     *bindings.DisputeGameFactory
	prevL2Block     uint64
	gerAddrL1       common.Address
	gerAddrL2       common.Address
}

func NewOP(cliCtx *cli.Context) (*OP, error) {
	l1ChainID := cliCtx.Int(l1ChainIDFlagName)
	l2ChainID := cliCtx.Int(l2ChainIDFlagName)
	l1RPC, l2RPC, err := goutils.LoadRPCs(l1ChainID, l2ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load RPCs: %w", err)
	}
	disputeGameFactoryProxyAddr := common.HexToAddress(cliCtx.String(gameFactoryFlagName))
	gameFactory, err := bindings.NewDisputeGameFactory(disputeGameFactoryProxyAddr, l1RPC)
	if err != nil {
		return nil, fmt.Errorf("failed to bind with disputeGameFactoryProxyAddr: %w", err)
	}

	return &OP{
		networkID:       uint64(cliCtx.Int(networkIDFlagName)),
		l1ChainID:       l1ChainID,
		l2ChainID:       l2ChainID,
		l1RPC:           l1RPC,
		l2RPC:           l2RPC,
		gameFactoryAddr: disputeGameFactoryProxyAddr,
		gameFactory:     gameFactory,
		prevL2Block:     uint64(cliCtx.Int(prevL2BlockNumFlagName)),
		gerAddrL1:       common.HexToAddress(cliCtx.String(gerL1FlagName)),
		gerAddrL2:       common.HexToAddress(cliCtx.String(gerL2FlagName)),
	}, nil
}

func (op *OP) GenerateChainProof(ctx context.Context) error {
	// Get current L1 block number
	l1BlockNum, err := op.l1RPC.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to call op.l1RPC.BlockNumber: %w", err)
	}

	// Get last game
	game, err := op.GetLastGame()
	if err != nil {
		return fmt.Errorf("failed to call op.GetLastGame: %w", err)
	}

	// GENERATE UNDERLAYING PROOFS IN PARALLEL
	chainSubProofsGroup, _ := errgroup.WithContext(ctx)

	// Generate consensus proof
	chainSubProofsGroup.Go(func() error {
		expectedOutputFile := path.Join("proof", fmt.Sprintf("chain%d", op.l2ChainID), fmt.Sprintf("game_%d.bin", game.index))
		if _, err := os.Stat(expectedOutputFile); err == nil {
			fmt.Println("consensus proof already exists (file): ", expectedOutputFile)
			return nil
		}
		proofCmd := fmt.Sprintf(`
		cargo run --release -- \
		--chain-id-l1 %d \
		--chain-id-l2 %d \
		--block-number-l1 %d \
		--prev-block-number-l2 %d \
		--new-block-number-l2 %d \
		--game-index %d \
		--game-factory-address %s \
		--claim-block-hash %s \
		--claim-state-root %s \
		--claim-message-passer-storage-root %s \
		--prove
	`,
			op.l1ChainID,        // --chain-id-l1
			op.l2ChainID,        // --chain-id-l2
			l1BlockNum,          // --block-number-l1
			op.prevL2Block,      // --prev-block-number-l2
			game.newL2Block,     // --new-block-number-l2
			game.index,          // --game-index
			op.gameFactoryAddr,  // --game-factory-address
			game.claimBlockHash, // --claim-block-hash
			game.claimStateRoot, // --claim-state-root
			game.claimMessage,   // --claim-message-passer-storage-root
		)
		fmt.Println("running the consensus proof command: ", proofCmd)
		cmd := exec.Command("bash", "-l", "-c", proofCmd)
		cmd.Dir = "consensus/host"
		msg, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to generate consensus proof:\n%s\n\n%w", string(msg), err)
		}
		fmt.Println("🚀🚀 🤝CONSENSUS PROOF GENERATED🤝 🚀🚀")
		return nil
	})
	chainSubProofsGroup.Go(func() error {
		return bridge.GenerateProof(
			ctx,
			op.l2RPC,
			"../bridge/host",
			uint64(op.l1ChainID),
			uint64(op.l2ChainID),
			l1BlockNum,
			op.prevL2Block,
			game.newL2Block,
			op.gerAddrL1,
			op.gerAddrL2,
		)
	})

	fmt.Println("generating proofs needed to build the chain proof...")
	if err = chainSubProofsGroup.Wait(); err != nil {
		return fmt.Errorf("failed to generate sub proofs: %w", err)
	}

	fmt.Println("generating chain proof")
	expectedOutputFile := path.Join("..", "chain-proofs", fmt.Sprintf("proof_chain_%d", op.networkID))
	if _, err := os.Stat(expectedOutputFile); err == nil {
		fmt.Println("consensus proof already exists (file): ", expectedOutputFile)
		return nil
	}
	proofCmd := fmt.Sprintf(`
		cargo run --release -- \
		--network-id %d \
		--chain-id %d \
		--prev-block-number-l2 %d \
		--new-block-number-l2 %d \
		--game-index %d \
		--prove
		`,
		op.networkID,    // --network-id
		op.l2ChainID,    // --chain-id
		op.prevL2Block,  // --prev-block-number-l2
		game.newL2Block, // --new-block-number-l2
		game.index,      // --game-index

	)
	fmt.Println("running the chain proof command: ", proofCmd)
	cmd := exec.Command("bash", "-l", "-c", proofCmd)
	cmd.Dir = "chain-proof/host"
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate chain proof:\n%s\n\n%w", string(msg), err)
	}

	fmt.Println("🚀🚀 🔗CHAIN PROOF GENERATED🔗 🚀🚀")
	return nil

}

type OPGame struct {
	newL2Block     uint64
	index          uint64
	claimBlockHash common.Hash
	claimStateRoot common.Hash
	claimMessage   common.Hash
}

func (o *OP) GetLastGame() (*OPGame, error) {
	// Get last game contract
	gameCounter, err := o.gameFactory.GameCount(nil)
	if err != nil {
		return nil, fmt.Errorf("o.gameFactory.GameCount: %w", err)
	}
	lastGameIndex := gameCounter.Sub(gameCounter, big.NewInt(1))
	gameCreation, err := o.gameFactory.GameAtIndex(nil, lastGameIndex)
	if err != nil {
		return nil, fmt.Errorf("o.gameFactory.GameAtIndex: %w", err)
	}
	gameContract, err := opfaultdisputegame.NewOpfulldisputegame(gameCreation.Proxy, o.l1RPC)
	if err != nil {
		return nil, fmt.Errorf("opfaultdisputegame.NewOpfulldisputegame: %w", err)
	}

	// check game output
	finalL2BlockNumber, err := gameContract.L2BlockNumber(nil)
	if err != nil {
		return nil, fmt.Errorf("gameContract.L2BlockNumber: %w", err)
	}
	finalRoot, err := gameContract.RootClaim(nil)
	if err != nil {
		return nil, fmt.Errorf("gameContract.RootClaim: %w", err)
	}
	finalOutputV0, err := OutputV0AtBlock(context.Background(), o.l2RPC, finalL2BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("OutputV0AtBlock final: %w", err)
	}
	expectedFinalRoot := eth.OutputRoot(finalOutputV0)
	if finalRoot != expectedFinalRoot {
		fmt.Printf("outputV0: %+v\n", finalOutputV0)
		fmt.Printf("final block number: %+v\n", finalL2BlockNumber.Int64())
		fmt.Printf("final claim root: %s\n", common.Hash(finalRoot).Hex())
		fmt.Printf("expectedFinalRoot: %s\n", common.Hash(expectedFinalRoot).Hex())
		return nil, errors.New("RootHash does not match")
	}

	// check game starting root
	initialL2Info, err := gameContract.StartingOutputRoot(nil)
	if err != nil {
		return nil, fmt.Errorf("gameContract.StartingOutputRoot: %w", err)
	}
	initialOutputV0, err := OutputV0AtBlock(context.Background(), o.l2RPC, initialL2Info.L2BlockNumber)
	if err != nil {
		return nil, fmt.Errorf("OutputV0AtBlock init: %w", err)
	}
	expectedInitialRoot := eth.OutputRoot(initialOutputV0)
	if initialL2Info.Root != expectedInitialRoot {
		fmt.Printf("initial outputV0: %+v\n", initialOutputV0)
		fmt.Printf("initial block number: %+v\n", initialL2Info.L2BlockNumber.Int64())
		fmt.Printf("initial claim root: %s\n", common.Hash(initialL2Info.Root).Hex())
		fmt.Printf("expectedInitialRoot: %s\n", common.Hash(expectedInitialRoot).Hex())
		return nil, errors.New("RootHash does not match")
	}

	fmt.Printf(`
========================== OP Game to be proved =================================
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
* Claim root: %s
=================================================================================
`,
		lastGameIndex, // game index
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

	return &OPGame{
		newL2Block:     finalL2BlockNumber.Uint64(),
		index:          lastGameIndex.Uint64(),
		claimBlockHash: finalOutputV0.BlockHash,
		claimStateRoot: common.Hash(finalOutputV0.StateRoot),
		claimMessage:   common.Hash(finalOutputV0.MessagePasserStorageRoot),
	}, nil
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
