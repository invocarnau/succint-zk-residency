package bridge

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/invocarnau/succint-zk-residency/goutils/contracts/gerl2"
)

var updateGEREvent = crypto.Keccak256Hash([]byte("InsertGlobalExitRoot(bytes32)"))

func SyncGERInjections(ctx context.Context, l2RPC *ethclient.Client, gerAddrL2 common.Address, fromBlock, toBlock uint64) ([]common.Hash, error) {
	sc, err := gerl2.NewGerl2(gerAddrL2, l2RPC)
	if err != nil {
		return nil, fmt.Errorf("gerl2.NewGerl2: %w", err)
	}
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		Addresses: []common.Address{gerAddrL2},
		ToBlock:   new(big.Int).SetUint64(toBlock),
	}
	logs, err := l2RPC.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("l2RPC.FilterLogs for GER: %w", err)
	}
	gers := []common.Hash{}
	for _, l := range logs {
		if l.Topics[0] == updateGEREvent {
			ger, err := sc.ParseInsertGlobalExitRoot(l)
			if err != nil {
				return nil, fmt.Errorf("sc.ParseInsertGlobalExitRoot: %w", err)
			}
			gers = append(gers, ger.NewGlobalExitRoot)
		}
	}

	return gers, nil
}

func GenerateProof(
	ctx context.Context,
	l2RPC *ethclient.Client,
	bridgeRelativePath string,
	l1ChainID, l2ChainID uint64,
	l1BlockNumber, prevL2BlockNumber, newL2BlockNumber uint64,
	gerAddrL1, gerAddrL2 common.Address,
) error {
	expectedOutputFile := path.Join(
		bridgeRelativePath, "..", "proof",
		fmt.Sprintf("chain%d", l2ChainID), fmt.Sprintf(
			"bridge_block_%d_to_%d_proof.bin", prevL2BlockNumber, newL2BlockNumber,
		),
	)
	if _, err := os.Stat(expectedOutputFile); err == nil {
		fmt.Println("bridge proof already exists (file): ", expectedOutputFile)
		return nil
	}
	fmt.Printf("synchronizing injected GERs from block %d to %d\n", prevL2BlockNumber, newL2BlockNumber)
	gers, err := SyncGERInjections(ctx, l2RPC, gerAddrL2, prevL2BlockNumber, newL2BlockNumber)
	if err != nil {
		return fmt.Errorf("syncGERInjections: %w", err)
	}
	gersStr := ""
	for _, ger := range gers {
		gersStr += ger.Hex() + ","
	}
	if len(gersStr) > 0 {
		gersStr = strings.TrimSuffix(gersStr, ",")
	}
	proofCmd := fmt.Sprintf(`
		cargo run --release -- \
		--chain-id-l1 %d \
		--chain-id-l2 %d \
		--block-number-l1 %d \
		--prev-block-number-l2 %d \
		--new-block-number-l2 %d \
		--contract-ger-l1 %s \
		--contract-ger-l2 %s \
		--imported-gers-hex "%s" \
		--prove
	`,
		l1ChainID,         // --chain-id-l1
		l2ChainID,         // --chain-id-l2
		l1BlockNumber,     // --block-number-l1
		prevL2BlockNumber, // --prev-block-number-l2
		newL2BlockNumber,  // --new-block-number-l2
		gerAddrL1.Hex(),   // --contract-ger-l1
		gerAddrL2.Hex(),   // --contract-ger-l2
		gersStr,           // --imported-gers-hex
	)
	fmt.Println("running the bridge proof command: ", proofCmd)
	cmd := exec.Command("bash", "-l", "-c", proofCmd)
	cmd.Dir = bridgeRelativePath
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate bridge proof:\n%s\n\n%w", string(msg), err)
	}
	fmt.Println("🚀🚀 📨BRIDGE PROOF GENERATED📨 🚀🚀")
	return nil
}
