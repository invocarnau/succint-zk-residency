package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/invocarnau/succint-zk-residency/bridge"
	"github.com/invocarnau/succint-zk-residency/goutils"
	"github.com/urfave/cli/v2"
)

const (
	proof                  = "proof"
	l1ChainIDFlagName      = "l1-chainid"
	l2ChainIDFlagName      = "l2-chainid"
	milestoneIdFlagName    = "milestone-id"
	milestoneHashFlagName  = "milestone-hash"
	prevL2BlockNumFlagName = "prev-l2-block-num"
	gerL1FlagName          = "ger-l1"
	gerL2FlagName          = "ger-l2"
)

func main() {
	app := cli.NewApp()
	app.Name = "PoS Aggregated Proof Generator"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "proof",
			Aliases: []string{},
			Usage:   "Generate an aggregated proof for PoS chain",
			Action:  generateProof,
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
				&cli.Uint64Flag{
					Name:     milestoneIdFlagName,
					Aliases:  []string{"id"},
					Usage:    "Milestone ID to generate consensus proof against",
					Required: false,
				},
				&cli.StringFlag{
					Name:     milestoneHashFlagName,
					Aliases:  []string{"hash"},
					Usage:    "Milestone hash to generate consensus proof against",
					Required: false,
				},
				&cli.Uint64Flag{
					Name:     prevL2BlockNumFlagName,
					Aliases:  []string{"from", "prev-block", "pb"},
					Usage:    "Specify the previous L2 block number, from which the proof will start, this should match the new block of the last proof",
					Required: false,
				},
				&cli.StringFlag{
					Name:     gerL1FlagName,
					Aliases:  []string{"gl1", "ger-addr-l1"},
					Usage:    "Address of the L1 GER contract",
					Required: false,
				},
				&cli.StringFlag{
					Name:     gerL2FlagName,
					Aliases:  []string{"gl2", "ger-addr-l2"},
					Usage:    "Address of the L2 GER contract",
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

func getGers(c *cli.Context) error {
	l1ChainID := c.Uint64(l1ChainIDFlagName)
	l2ChainID := c.Uint64(l2ChainIDFlagName)
	gerL2 := c.String(gerL2FlagName)
	_, l2Rpc, err := goutils.LoadRPCs(int(l1ChainID), int(l2ChainID))
	if err != nil {
		return err
	}
	gers, err := bridge.SyncGERInjections(context.Background(), l2Rpc, common.HexToAddress(gerL2), 13296845, 13297750)
	if err != nil {
		return err
	}
	fmt.Println("gers:", gers)
	gersStr := ""
	for _, ger := range gers {
		gersStr += ger.Hex() + ","
	}
	if len(gersStr) > 0 {
		gersStr = strings.TrimSuffix(gersStr, ",")
	}
	fmt.Println("gersStr:", gersStr)
	return nil
}

func generateProof(c *cli.Context) error {
	l1ChainID := c.Uint64(l1ChainIDFlagName)
	l2ChainID := c.Uint64(l2ChainIDFlagName)
	milestoneId := c.Uint64(milestoneIdFlagName)
	milestoneHash := c.String(milestoneHashFlagName)
	prevL2Block := c.Uint64(prevL2BlockNumFlagName)
	gerL1 := c.String(gerL1FlagName)
	gerL2 := c.String(gerL2FlagName)

	l1Rpc, l2Rpc, err := goutils.LoadRPCs(int(l1ChainID), int(l2ChainID))
	if err != nil {
		return err
	}
	l1BlockNumber, err := l1Rpc.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	l2Block, err := findL2Block(uint64(milestoneId))
	if err != nil {
		return err
	}

	// Generate both proofs in parallel
	var wg sync.WaitGroup
	var err1, err2 bool
	go func() {
		wg.Add(1)
		err := generateConsensusProof(l1ChainID, milestoneId, milestoneHash, l1BlockNumber, prevL2Block, l2Block)
		if err != nil {
			fmt.Println("failed to generate consensus proofs:", err)
			err1 = true
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		err := bridge.GenerateProof(context.Background(), l2Rpc, "", l1ChainID, l2ChainID, l1BlockNumber, prevL2Block, l2Block, common.HexToAddress(gerL1), common.HexToAddress(gerL2))
		if err != nil {
			fmt.Println("failed to generate bridge proofs:", err)
			err2 = true
		}
		wg.Done()
	}()
	wg.Wait()
	if err1 || err2 {
		return fmt.Errorf("failed to generate proofs, consensus proof result: %v, bridge proof result: %v", err1, err2)
	}

	// Generate aggregated chain proof
	err = generateChainProof(l2ChainID, prevL2Block, l2Block)
	if err != nil {
		return err
	}

	return nil
}

func generateConsensusProof(l2ChainId uint64, milestoneId uint64, milestoneHash string, l1BlockNumber uint64, prevL2BlockNumber uint64, newL2BlockNumber uint64) error {
	expectedOutputFile := path.Join("proofs", fmt.Sprintf("chain%d/consensus_block_%d_to_%d", l2ChainId, prevL2BlockNumber, newL2BlockNumber))
	if _, err := os.Stat(expectedOutputFile); err == nil {
		fmt.Println("consensus proof already exists (file): ", expectedOutputFile)
		return nil
	}

	// cargo run --bin operator --release -- --milestone-id 332570 \
	// --milestone-hash 0xa25a2d394de5d8bede49a90c870fcdae72cdcf5e7dd26c117e17c3ffa9d35ec7 \
	// --l1-block-number 6894731 \
	// --prev-l2-block-number 13296845 --new-l2-block-number 13298032 \
	// --prove
	// Generate the consensus proof command for operator
	proofCmd := fmt.Sprintf(
		`cargo run --bin operator --release -- --milestone-id %d \
		--milestone-hash %s \
		--l1-block-number %d \
		--prev-l2-block-number %d \
		--new-l2-block-number %d \
		--prove`,
		milestoneId,
		milestoneHash,
		l1BlockNumber,
		prevL2BlockNumber,
		newL2BlockNumber,
	)
	fmt.Println("running the consensus proof command:", proofCmd)
	cmd := exec.Command("bash", "-l", "-c", proofCmd)
	cmd.Dir = "consensus/host"
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate consensus proof:\n%s\n\n%w", string(msg), err)
	}
	fmt.Println("🚀🚀 🤝CONSENSUS PROOF GENERATED🤝 🚀🚀")
	return nil
}

func generateChainProof(l2ChainId uint64, prevL2BlockNumber uint64, newL2BlockNumber uint64) error {
	// cargo run --release --bin main -- --prev-l2-block-number 13296845 --new-l2-block-number 13298032 --prove
	proofCmd := fmt.Sprintf(
		`cargo run --bin main --release -- --prev-l2-block-number %d --new-l2-block-number %d --prove`,
		prevL2BlockNumber,
		newL2BlockNumber,
	)
	fmt.Println("running the chain proof command:", proofCmd)
	cmd := exec.Command("bash", "-l", "-c", proofCmd)
	cmd.Dir = "chain-proof/host"
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate chain proof for pos:\n%s\n\n%w", string(msg), err)
	}
	fmt.Println("🚀🚀 🤝CHAIN PROOF FOR POS GENERATED🤝 🚀🚀")
	return nil
}

func findL2Block(id uint64) (uint64, error) {
	heimdallEndpoint, err := goutils.LoadHeimdallEndpoint()
	if err != nil {
		return 0, err
	}

	url := fmt.Sprintf("%s/milestone/%d", heimdallEndpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	type Response struct {
		Result struct {
			EndBlock uint64 `json:"end_block"`
		} `json:"result"`
	}

	// Unmarshal the JSON response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	return response.Result.EndBlock, nil
}
