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
					Required: true,
				},
				&cli.StringFlag{
					Name:     milestoneHashFlagName,
					Aliases:  []string{"hash"},
					Usage:    "Milestone hash to generate consensus proof against",
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func generateProof(c *cli.Context) error {
	l1ChainID := c.Uint64(l1ChainIDFlagName)
	l2ChainID := c.Uint64(l2ChainIDFlagName)
	milestoneId := c.Uint64(milestoneIdFlagName)
	milestoneHash := c.String(milestoneHashFlagName)
	prevL2Block := c.Int(prevL2BlockNumFlagName)
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
		err := GenerateConsensusProof(milestoneId, milestoneHash, l1BlockNumber)
		if err != nil {
			fmt.Println("failed to generate consensus proofs:", err)
			err1 = true
		}
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		err := bridge.GenerateProof(context.Background(), l2Rpc, "", l1ChainID, l2ChainID, l1BlockNumber, uint64(prevL2Block), l2Block, common.HexToAddress(gerL1), common.HexToAddress(gerL2))
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

	// TODO: Generate chain proof if everything goes well

	return nil
}

func GenerateConsensusProof(milestoneId uint64, milestoneHash string, l1BlockNumber uint64) error {
	expectedOutputFile := path.Join("proofs/consensus", fmt.Sprintf("milestone_%d", milestoneId))
	if _, err := os.Stat(expectedOutputFile); err == nil {
		fmt.Println("consensus proof already exists (file): ", expectedOutputFile)
		return nil
	}

	// cargo run --bin operator --release -- --milestone-id 329745 --milestone-hash 0x708d99ad07a0cee3e696197e838073000c27fde2daa8a63290dff27c5e4932b4 --prove
	// Generate the consensus proof command for operator
	proofCmd := fmt.Sprintf(
		`cargo run --bin operator --release \
		--milestone-id %d \
		--milestone-hash %s \
		--l1_block_number %d \
		--prove`,
		milestoneId,
		milestoneHash,
		l1BlockNumber,
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
