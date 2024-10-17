package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/etrog/polygonzkevmbridgev2"
	"github.com/0xPolygon/cdk/bridgesync"
	"github.com/0xPolygon/cdk/claimsponsor"
	"github.com/0xPolygon/cdk/l1infotreesync"
	cdkClient "github.com/0xPolygon/cdk/rpc/client"
	"github.com/ethereum/go-ethereum/common"
	"github.com/invocarnau/succint-zk-residency/goutils"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

const (
	originTokenAddrFlagName             = "origin-token-addr"
	originNetworkFlagName               = "origin-network"
	destinationNetworkFlagName          = "dest-net"
	destinationAddressFlagName          = "dest-addr"
	amountFlagName                      = "amount"
	depositCountFlagName                = "deposit-count"
	originBridgeServiceURLFlagName      = "origin-cdk"
	destinationBridgeServiceURLFlagName = "dest-cdk"
	metadataFlagName                    = "metadata"
	privateKeyFlagName                  = "pk"
)

var claimFlags []cli.Flag = []cli.Flag{
	&cli.StringFlag{
		Name:     originTokenAddrFlagName,
		Aliases:  []string{"ota"},
		Usage:    "Origin token address of the bridge",
		Required: true,
	},
	&cli.IntFlag{
		Name:     originNetworkFlagName,
		Aliases:  []string{"on"},
		Usage:    "Origin network ID",
		Required: true,
	},
	&cli.IntFlag{
		Name:     destinationNetworkFlagName,
		Aliases:  []string{"dn"},
		Usage:    "Destination network ID",
		Required: true,
	},
	&cli.StringFlag{
		Name:     destinationAddressFlagName,
		Aliases:  []string{"da"},
		Usage:    "Destination address of the bridge",
		Required: true,
	},
	&cli.StringFlag{
		Name:     amountFlagName,
		Aliases:  []string{"a"},
		Usage:    "Amount of the bridge",
		Required: true,
	},
	&cli.IntFlag{
		Name:     depositCountFlagName,
		Aliases:  []string{"dc"},
		Usage:    "Deposit count of the bridge",
		Required: true,
	},
	&cli.StringFlag{
		Name:     originBridgeServiceURLFlagName,
		Aliases:  []string{"orig-cdk-url"},
		Usage:    "URL of the CDK service where the bridge has originated",
		Required: true,
	},
	&cli.StringFlag{
		Name:     destinationBridgeServiceURLFlagName,
		Aliases:  []string{"dest-cdk-url"},
		Usage:    "URL of the CDK service where the bridge is headed",
		Required: true,
	},
	&cli.StringFlag{
		Name:     metadataFlagName,
		Aliases:  []string{"meta"},
		Usage:    "Metadata of the bridge",
		Required: true,
	},
	&cli.StringFlag{
		Name:     privateKeyFlagName,
		Aliases:  []string{},
		Usage:    "Private key to send the tx. Only for claim (not claim sponsor)",
		Required: false,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "OP getter"
	app.Version = "v0.0.1"
	app.Commands = []*cli.Command{
		{
			Name:    "build-config",
			Aliases: []string{},
			Usage:   "Take the CDK template config and builds it for FEP, OP and PoS",
			Action:  buildConfigCDK,
			Flags:   []cli.Flag{},
		},
		{
			Name:    "sponsor-claim",
			Aliases: []string{},
			Usage:   "Send a sponsor claim request",
			Action:  sponsorClaim,
			Flags:   claimFlags,
		},
		{
			Name:    "claim",
			Aliases: []string{},
			Usage:   "Send a claim request",
			Action:  claim,
			Flags:   claimFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func buildConfigCDK(cli *cli.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	for _, chain := range []string{"op", "fep", "pos"} {
		if err := buildConfigFromTemplate(chain); err != nil {
			return err
		}
	}

	return nil
}

func sponsorClaim(cli *cli.Context) error {
	// PARSE INPUTS
	bridgeClientDest := cdkClient.NewClient(cli.String(destinationBridgeServiceURLFlagName))
	bridgeClientOrig := cdkClient.NewClient(cli.String(originBridgeServiceURLFlagName))
	originNetwork := cli.Int(originNetworkFlagName)
	destinationNetwork := cli.Int(destinationNetworkFlagName)
	originTokenAddress := common.HexToAddress(cli.String(originTokenAddrFlagName))
	destinationAddress := common.HexToAddress(cli.String(destinationAddressFlagName))
	amount, _ := big.NewInt(0).SetString(cli.String(amountFlagName), 10)
	depositCount := uint32(cli.Int(depositCountFlagName))
	claim := claimsponsor.Claim{
		LeafType:           0,
		OriginNetwork:      uint32(originNetwork),
		OriginTokenAddress: originTokenAddress,
		DestinationNetwork: uint32(destinationNetwork),
		DestinationAddress: destinationAddress,
		Amount:             amount,
		Metadata:           common.Hex2Bytes(cli.String(metadataFlagName)),
	}
	fmt.Printf("claim: %+v\n", claim)

	// GET L1 info tree index in which bridge included
	var err error
	var bridgeIncluddedAtIndex uint32
	for i := 0; i < 40; i++ { // block needs to be finalised, takes ~32s
		fmt.Println("bridgeClientOrig.L1InfoTreeIndexForBridge")
		bridgeIncluddedAtIndex, err = bridgeClientOrig.L1InfoTreeIndexForBridge(claim.OriginNetwork, depositCount)
		if err == nil {
			break
		} else {
			fmt.Println("error! ", err)
		}
		time.Sleep(time.Second * 2)
	}
	fmt.Println("Bridge includded at L1 Info Tree Index: ", bridgeIncluddedAtIndex)

	// GET L1 INFO TREE INDEX
	fmt.Println("getting info already injected on L2")
	var info *l1infotreesync.L1InfoTreeLeaf
	for {
		fmt.Println("bridgeClientDest.InjectedInfoAfterIndex")
		info, err = bridgeClientDest.InjectedInfoAfterIndex(claim.DestinationNetwork, bridgeIncluddedAtIndex)
		if err == nil {
			break
		} else {
			fmt.Println("err! ", err)
		}
		time.Sleep(time.Second * 2)
	}

	// GET PROOF
	fmt.Printf("Info associated to the first GER injected on L2 after index %d: %+v\n", bridgeIncluddedAtIndex, info)
	proof, err := bridgeClientOrig.ClaimProof(0, depositCount, info.L1InfoTreeIndex)
	if err != nil {
		return fmt.Errorf("err ClaimProof: %w", err)
	}
	fmt.Printf("ClaimProof received from bridge service\n")

	fmt.Println("Sending claim tx")
	claim.ProofLocalExitRoot = proof.ProofLocalExitRoot
	claim.ProofRollupExitRoot = proof.ProofRollupExitRoot
	claim.GlobalIndex = bridgesync.GenerateGlobalIndex(true, claim.DestinationNetwork-1, depositCount)
	claim.MainnetExitRoot = info.MainnetExitRoot
	claim.RollupExitRoot = info.RollupExitRoot

	// SPONSOR CLAIM
	fmt.Println("waiting for service to send claim on behalf of the user...")
	if err := bridgeClientDest.SponsorClaim(claim); err != nil {
		return fmt.Errorf("failed to request claim sponsor: %w", err)
	}
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 2)
		status, err := bridgeClientDest.GetSponsoredClaimStatus(claim.GlobalIndex)
		fmt.Println("sponsored claim status: ", status)
		if err != nil {
			fmt.Println("error getting sponsored claim status: ", err)
			continue
		}
		if claimsponsor.FailedClaimStatus == status {
			return fmt.Errorf("claim response unexpected: %s", status)
		}
		if status == claimsponsor.SuccessClaimStatus {
			break
		}
	}
	fmt.Println("claim request succeed")
	return nil
}

func claim(cli *cli.Context) error {
	// PARSE INPUTS
	bridgeClientDest := cdkClient.NewClient(cli.String(destinationBridgeServiceURLFlagName))
	bridgeClientOrig := cdkClient.NewClient(cli.String(originBridgeServiceURLFlagName))
	originNetwork := cli.Int(originNetworkFlagName)
	destinationNetwork := cli.Int(destinationNetworkFlagName)
	originTokenAddress := common.HexToAddress(cli.String(originTokenAddrFlagName))
	destinationAddress := common.HexToAddress(cli.String(destinationAddressFlagName))
	amount, _ := big.NewInt(0).SetString(cli.String(amountFlagName), 10)
	depositCount := uint32(cli.Int(depositCountFlagName))
	claim := claimsponsor.Claim{
		LeafType:           0,
		OriginNetwork:      uint32(originNetwork),
		OriginTokenAddress: originTokenAddress,
		DestinationNetwork: uint32(destinationNetwork),
		DestinationAddress: destinationAddress,
		Amount:             amount,
		Metadata:           common.Hex2Bytes(cli.String(metadataFlagName)),
	}
	fmt.Printf("claim: %+v\n", claim)

	// GET L1 info tree index in which bridge included
	var err error
	var bridgeIncluddedAtIndex uint32
	for i := 0; i < 40; i++ { // block needs to be finalised, takes ~32s
		fmt.Println("bridgeClientOrig.L1InfoTreeIndexForBridge")
		bridgeIncluddedAtIndex, err = bridgeClientOrig.L1InfoTreeIndexForBridge(claim.OriginNetwork, depositCount)
		if err == nil {
			break
		} else {
			fmt.Println("error! ", err)
		}
		time.Sleep(time.Second * 2)
	}
	fmt.Println("Bridge includded at L1 Info Tree Index: ", bridgeIncluddedAtIndex)

	// GET L1 INFO TREE INDEX
	fmt.Println("getting info already injected on L2")
	var info *l1infotreesync.L1InfoTreeLeaf
	for {
		fmt.Println("bridgeClientDest.InjectedInfoAfterIndex")
		info, err = bridgeClientDest.InjectedInfoAfterIndex(claim.DestinationNetwork, bridgeIncluddedAtIndex)
		if err == nil {
			break
		} else {
			fmt.Println("err! ", err)
		}
		time.Sleep(time.Second * 2)
	}

	// GET PROOF
	fmt.Printf("Info associated to the first GER injected on L2 after index %d: %+v\n", bridgeIncluddedAtIndex, info)
	proof, err := bridgeClientOrig.ClaimProof(0, depositCount, info.L1InfoTreeIndex)
	if err != nil {
		return fmt.Errorf("err ClaimProof: %w", err)
	}
	fmt.Printf("ClaimProof received from bridge service\n")

	fmt.Println("Requesting service to sponsor claim")
	claim.ProofLocalExitRoot = proof.ProofLocalExitRoot
	claim.ProofRollupExitRoot = proof.ProofRollupExitRoot
	claim.GlobalIndex = bridgesync.GenerateGlobalIndex(true, claim.DestinationNetwork-1, depositCount)
	claim.MainnetExitRoot = info.MainnetExitRoot
	claim.RollupExitRoot = info.RollupExitRoot

	// SPONSOR CLAIM
	fmt.Println("waiting for service to send claim on behalf of the user...")
	abi, err := polygonzkevmbridgev2.Polygonzkevmbridgev2MetaData.GetAbi()
	if err != nil {
		return err
	}
	data, err := abi.Pack(
		"claimAsset",
		claim.ProofLocalExitRoot,  // bytes32[32] smtProofLocalExitRoot
		claim.ProofRollupExitRoot, // bytes32[32] smtProofRollupExitRoot
		claim.GlobalIndex,         // uint256 globalIndex
		claim.MainnetExitRoot,     // bytes32 mainnetExitRoot
		claim.RollupExitRoot,      // bytes32 rollupExitRoot
		claim.OriginNetwork,       // uint32 originNetwork
		claim.OriginTokenAddress,  // address originTokenAddress,
		claim.DestinationNetwork,  // uint32 destinationNetwork
		claim.DestinationAddress,  // address destinationAddress
		claim.Amount,              // uint256 amount
		claim.Metadata,            // bytes metadata
	)
	if err != nil {
		return err
	}

	fmt.Printf("Data to send the claim tx: %s", common.Bytes2Hex(data))
	return nil
}

func buildConfigFromTemplate(l2Separator string) error {
	file, err := os.ReadFile("./config/cdk-template.toml")
	if err != nil {
		return err
	}
	l1ChainIDStr := os.Getenv("XXX_CHAINIDL1_XXX")
	l1ChainID, err := strconv.Atoi(l1ChainIDStr)
	if err != nil {
		return err
	}
	l2ChainIDStr := os.Getenv("XXX_CHAINIDL2_XXX_" + l2Separator)
	l2ChainID, err := strconv.Atoi(l2ChainIDStr)
	if err != nil {
		return err
	}
	l1URL, l2URL, err := goutils.LoadRPCURLs(l1ChainID, l2ChainID)
	if err != nil {
		return err
	}
	updatedConfig := strings.ReplaceAll(string(file), "XXX_GERL1_XXX", os.Getenv("XXX_GERL1_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_ROLLUPMAN_XXX", os.Getenv("XXX_ROLLUPMAN_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_INITBLOCK_L1_XXX", os.Getenv("XXX_INITBLOCK_L1_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_AGGORACLESENDERADDR_XXX", os.Getenv("XXX_AGGORACLESENDERADDR_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_AGGORACLEKEYPASS_XXX", os.Getenv("XXX_AGGORACLEKEYPASS_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_CLAIMSPONSORSENDERADDR_XXX", os.Getenv("XXX_CLAIMSPONSORSENDERADDR_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_CLAIMSPONSORKEYPASS_XXX", os.Getenv("XXX_CLAIMSPONSORKEYPASS_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_BRIDGEL1_XXX", os.Getenv("XXX_BRIDGEL1_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_CHAINIDL1_XXX", os.Getenv("XXX_CHAINIDL1_XXX"))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_RPCL1_XXX", l1URL)
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_RPCL2_XXX", l2URL)

	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_INITBLOCK_L2_XXX", os.Getenv("XXX_INITBLOCK_L2_XXX_"+l2Separator))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_GERL2_XXX", os.Getenv("XXX_GERL2_XXX_"+l2Separator))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_CHAINIDL2_XXX", os.Getenv("XXX_CHAINIDL2_XXX_"+l2Separator))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_BRIDGEL2_XXX", os.Getenv("XXX_BRIDGEL2_XXX_"+l2Separator))
	updatedConfig = strings.ReplaceAll(updatedConfig, "XXX_ROLLUPADDR_XXX", os.Getenv("XXX_ROLLUPADDR_XXX_"+l2Separator))

	return os.WriteFile(fmt.Sprintf("./config/cdk-%s.toml", l2Separator), []byte(updatedConfig), 0644)
}
