package testaggoracle

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmbridgev2"
	gerContractL1 "github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmglobalexitrootv2"
	"github.com/0xPolygon/cdk/bridgesync"
	"github.com/0xPolygon/cdk/claimsponsor"
	"github.com/0xPolygon/cdk/l1infotreesync"
	cdkClient "github.com/0xPolygon/cdk/rpc/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/invocarnau/succint-zk-residency/fep-type-1/e2etest/gerl2"
	"github.com/invocarnau/succint-zk-residency/fep-type-1/e2etest/transparentupgradableproxy"
	"github.com/stretchr/testify/require"
)

const (
	l1ChainID       = 1337
	l2ChainID       = 42069
	l1URL           = "http://localhost:8545"
	l2URL           = "http://localhost:8555"
	alreadyDeployed = true
)

var (
	gerAddrL2AlreadyDeployed    = common.HexToAddress("0x8058D80131e6F57E99830Dce403BBAF4e64C9b8A")
	bridgeAddrL2AlreadyDeployed = common.HexToAddress("0xb0a5546A0Efd8950D8964a9dB66DFF5569EEfDE7")
)

func TestBridgeEVM(t *testing.T) {
	// defer exec.Command("bash", "-l", "-c", "docker compose down").Run()
	authL1, authL2 := loadAuth(t)
	fmt.Println("running L1 network (turning up docker container)...")
	clientL1, _, gerL1, _, bridgeL1 := runL1(t)
	fmt.Println("running L2 network (turning up docker container + deploy contracts)...")
	gerAddrL2, _, bridgeAddrL2, bridgeL2 := runL2(t, authL2)
	if !alreadyDeployed {
		fmt.Println("running CDK client for L2 (turning up docker container)...")
		editConfig(t, gerAddrL2, bridgeAddrL2)
		runCDK(t)
	}
	runBridgeL1toL2Test(t, clientL1, authL1, authL2, gerL1, bridgeL1, bridgeL2)
}

func loadAuth(t *testing.T) (*bind.TransactOpts, *bind.TransactOpts) {
	keystoreEncrypted, err := os.ReadFile("./config/aggoracle.keystore")
	require.NoError(t, err)
	key, err := keystore.DecryptKey(keystoreEncrypted, "testonly")
	require.NoError(t, err)
	authL1, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(l1ChainID))
	require.NoError(t, err)
	authL2, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, new(big.Int).SetUint64(l2ChainID))
	require.NoError(t, err)
	return authL1, authL2
}

func runL1(t *testing.T) (
	*ethclient.Client,
	common.Address,
	*gerContractL1.Polygonzkevmglobalexitrootv2,
	common.Address,
	*polygonzkevmbridgev2.Polygonzkevmbridgev2,
) {
	gerAddr := common.HexToAddress("0x8A791620dd6260079BF849Dc5567aDC3F2FdC318")
	bridgeAddr := common.HexToAddress(("0xFe12ABaa190Ef0c8638Ee0ba9F828BF41368Ca0E"))
	if alreadyDeployed {
		msg, err := exec.Command("bash", "-l", "-c", "docker compose up -d test-fep-type1-l1").CombinedOutput()
		require.NoError(t, err, string(msg))
		time.Sleep(time.Second * 2)
	}
	client, err := ethclient.Dial(l1URL)
	require.NoError(t, err)
	gerContract, err := gerContractL1.NewPolygonzkevmglobalexitrootv2(gerAddr, client)
	require.NoError(t, err)
	bridgeContract, err := polygonzkevmbridgev2.NewPolygonzkevmbridgev2(bridgeAddr, client)
	require.NoError(t, err)
	return client, gerAddr, gerContract, bridgeAddr, bridgeContract
}

func runL2(t *testing.T, auth *bind.TransactOpts) (
	common.Address,
	*gerl2.Gerl2,
	common.Address,
	*polygonzkevmbridgev2.Polygonzkevmbridgev2,
) {
	client, err := ethclient.Dial(l2URL)
	require.NoError(t, err)
	if !alreadyDeployed {
		msg, err := exec.Command("bash", "-l", "-c", "docker compose up -d test-fep-type1-l2").CombinedOutput()
		require.NoError(t, err, string(msg))
		time.Sleep(time.Second * 2)
		require.NoError(t, err)

		// create tmp auth to deploy contracts
		ctx := context.Background()
		privateKeyL2, err := crypto.GenerateKey()
		require.NoError(t, err)
		authDeployer, err := bind.NewKeyedTransactorWithChainID(privateKeyL2, big.NewInt(l2ChainID))
		require.NoError(t, err)

		// fund deployer
		nonce, err := client.PendingNonceAt(ctx, auth.From)
		require.NoError(t, err)
		amountToTransfer, _ := new(big.Int).SetString("1000000000000000000", 10) //nolint:gomnd
		gasPrice, err := client.SuggestGasPrice(ctx)
		require.NoError(t, err)
		gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{From: auth.From, To: &authDeployer.From, Value: amountToTransfer})
		require.NoError(t, err)
		tx := types.NewTransaction(nonce, authDeployer.From, amountToTransfer, gasLimit, gasPrice, nil)
		signedTx, err := auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)
		time.Sleep(time.Second * 2)
		balance, err := client.BalanceAt(ctx, authDeployer.From, nil)
		require.NoError(t, err)
		require.Equal(t, amountToTransfer, balance)

		// fund bridge
		precalculatedBridgeAddr := crypto.CreateAddress(authDeployer.From, 1)
		tx = types.NewTransaction(nonce+1, precalculatedBridgeAddr, amountToTransfer, gasLimit, gasPrice, nil)
		signedTx, err = auth.Signer(auth.From, tx)
		require.NoError(t, err)
		err = client.SendTransaction(ctx, signedTx)
		require.NoError(t, err)
		time.Sleep(time.Second * 2)
		balance, err = client.BalanceAt(ctx, precalculatedBridgeAddr, nil)
		require.NoError(t, err)
		require.Equal(t, amountToTransfer, balance)

		// deploy bridge impl
		bridgeImplementationAddr, _, _, err := polygonzkevmbridgev2.DeployPolygonzkevmbridgev2(authDeployer, client)
		require.NoError(t, err)
		time.Sleep(time.Second * 2)

		// deploy bridge proxy
		nonce, err = client.PendingNonceAt(ctx, authDeployer.From)
		require.NoError(t, err)
		precalculatedAddr := crypto.CreateAddress(authDeployer.From, nonce+1)
		bridgeABI, err := polygonzkevmbridgev2.Polygonzkevmbridgev2MetaData.GetAbi()
		require.NoError(t, err)
		if bridgeABI == nil {
			err = errors.New("GetABI returned nil")
			require.NoError(t, err)
		}
		dataCallProxy, err := bridgeABI.Pack("initialize",
			uint32(1),        //network ID
			common.Address{}, // gasTokenAddressMainnet"
			uint32(0),        // gasTokenNetworkMainnet
			precalculatedAddr,
			common.Address{},
			[]byte{}, // gasTokenMetadata
		)
		require.NoError(t, err)
		code, err := client.CodeAt(ctx, bridgeImplementationAddr, nil)
		require.NoError(t, err)
		require.NotEqual(t, len(code), 0)
		bridgeAddr, _, _, err := transparentupgradableproxy.DeployTransparentupgradableproxy(
			authDeployer,
			client,
			bridgeImplementationAddr,
			authDeployer.From,
			dataCallProxy,
		)
		require.NoError(t, err)
		if bridgeAddr != precalculatedBridgeAddr {
			err = fmt.Errorf("error calculating bridge addr. Expected: %s. Actual: %s", precalculatedBridgeAddr, bridgeAddr)
			require.NoError(t, err)
		}
		time.Sleep(time.Second * 2)
		bridgeContract, err := polygonzkevmbridgev2.NewPolygonzkevmbridgev2(bridgeAddr, client)
		require.NoError(t, err)
		checkGERAddr, err := bridgeContract.GlobalExitRootManager(&bind.CallOpts{})
		require.NoError(t, err)
		if precalculatedAddr != checkGERAddr {
			err = errors.New("error deploying bridge")
			require.NoError(t, err)
		}

		// deploy GER
		gerAddr, _, gerContract, err := gerl2.DeployGerl2(authDeployer, client, auth.From)
		require.NoError(t, err)
		time.Sleep(time.Second * 2)
		fmt.Println("gerAddr ", gerAddr)
		fmt.Println("bridgeAddr ", bridgeAddr)

		return gerAddr, gerContract, bridgeAddr, bridgeContract
	} else {
		gerContract, err := gerl2.NewGerl2(gerAddrL2AlreadyDeployed, client)
		require.NoError(t, err)
		bridgeContract, err := polygonzkevmbridgev2.NewPolygonzkevmbridgev2(bridgeAddrL2AlreadyDeployed, client)
		require.NoError(t, err)
		return gerAddrL2AlreadyDeployed, gerContract, bridgeAddrL2AlreadyDeployed, bridgeContract
	}
}

func editConfig(t *testing.T, gerL2, bridgeL2 common.Address) {
	file, err := os.ReadFile("./config/template_cdk.toml")
	require.NoError(t, err)
	updatedConfig := strings.Replace(string(file), "XXX_GlobalExitRootL2", gerL2.Hex(), 2)
	updatedConfig = strings.Replace(updatedConfig, "XXX_BridgeL2", bridgeL2.Hex(), 2)
	err = os.WriteFile("./config/cdk.toml", []byte(updatedConfig), 0644)
	require.NoError(t, err)
}

func runCDK(t *testing.T) {
	msg, err := exec.Command("bash", "-l", "-c", "docker compose up -d test-fep-type1-cdk").CombinedOutput()
	require.NoError(t, err, string(msg))
	time.Sleep(time.Second * 2)
	require.NoError(t, err)
}

func runBridgeL1toL2Test(
	t *testing.T,
	clientL1 *ethclient.Client,
	authL1 *bind.TransactOpts,
	authL2 *bind.TransactOpts,
	gerL1Contract *gerContractL1.Polygonzkevmglobalexitrootv2,
	bridgeL1 *polygonzkevmbridgev2.Polygonzkevmbridgev2,
	bridgeL2 *polygonzkevmbridgev2.Polygonzkevmbridgev2,
) {
	l2NetworkID := uint32(1)
	bridgeClient := cdkClient.NewClient("http://localhost:5576")
	for i := 0; i < 1000; i++ {
		// Send bridge L1 -> L2
		fmt.Println("--- ITERATION ", i)
		fmt.Println("sending bridge tx to L1")
		amount := big.NewInt(int64(i + 1))
		authL1.Value = amount
		claim := claimsponsor.Claim{
			LeafType:           0,
			OriginNetwork:      0,
			OriginTokenAddress: common.Address{},
			DestinationNetwork: l2NetworkID,
			DestinationAddress: authL2.From,
			Amount:             amount,
			Metadata:           nil,
		}
		gerBefore, err := gerL1Contract.GetLastGlobalExitRoot(nil)
		require.NoError(t, err)
		tx, err := bridgeL1.BridgeAsset(authL1, claim.DestinationNetwork, claim.DestinationAddress, claim.Amount, claim.OriginTokenAddress, true, nil)
		require.NoError(t, err)
		time.Sleep(time.Second * 2)
		gerAfter, err := gerL1Contract.GetLastGlobalExitRoot(nil)
		require.NoError(t, err)
		require.NotEqual(t, gerBefore, gerAfter)
		fmt.Println("bridge sent")

		// Interact with bridge service
		fmt.Println("interacting with bridges service:")
		fmt.Println("waiting for the bridge to be finalised")
		receipt, err := clientL1.TransactionReceipt(context.TODO(), tx.Hash())
		require.NoError(t, err)
		bridgeEvent, err := bridgeL1.ParseBridgeEvent(*receipt.Logs[0])
		require.NoError(t, err)
		require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)
		depositCount := bridgeEvent.DepositCount
		var bridgeIncluddedAtIndex uint32
		found := false
		for i := 0; i < 40; i++ { // block needs to be finalised, takes ~32s
			bridgeIncluddedAtIndex, err = bridgeClient.L1InfoTreeIndexForBridge(0, depositCount)
			if err == nil {
				found = true
				break
			}
			time.Sleep(time.Second * 2)
		}
		require.True(t, found)
		fmt.Println("Bridge includded at L1 Info Tree Index: ", bridgeIncluddedAtIndex)
		var info *l1infotreesync.L1InfoTreeLeaf
		found = false
		for i := 0; i < 34; i++ {
			info, err = bridgeClient.InjectedInfoAfterIndex(l2NetworkID, bridgeIncluddedAtIndex)
			if err == nil {
				found = true
				break
			}
			time.Sleep(time.Second * 2)
		}
		require.True(t, found)
		require.NoError(t, err)
		fmt.Printf("Info associated to the first GER injected on L2 after index %d: %+v\n", bridgeIncluddedAtIndex, info)
		proof, err := bridgeClient.ClaimProof(0, depositCount, info.L1InfoTreeIndex)
		require.NoError(t, err)
		fmt.Printf("ClaimProof received from bridge service\n")
		claim.ProofLocalExitRoot = proof.ProofLocalExitRoot
		claim.ProofRollupExitRoot = proof.ProofRollupExitRoot
		claim.GlobalIndex = bridgesync.GenerateGlobalIndex(true, claim.DestinationNetwork-1, depositCount)
		claim.MainnetExitRoot = info.MainnetExitRoot
		claim.RollupExitRoot = info.RollupExitRoot
		err = bridgeClient.SponsorClaim(claim)
		require.NoError(t, err)
		fmt.Println("Requesting service to sponsor claim")
		fmt.Println("waiting for service to send claim on behalf of the user...")
		found = false
		for i := 0; i < 20; i++ {
			status, err := bridgeClient.GetSponsoredClaimStatus(claim.GlobalIndex)
			fmt.Println("sponsored claim status: ", status)
			require.NoError(t, err)
			require.NotEqual(t, claimsponsor.FailedClaimStatus, status)
			if status == claimsponsor.SuccessClaimStatus {
				found = true
				break
			}
			time.Sleep(time.Second * 2)
		}
		require.True(t, found)
		fmt.Println("service reports that the claim tx is succesful")

		// check that the bridge is claimed on L2
		fmt.Println("checking if bridge is claimed on L2...")
		isClaimed, err := bridgeL2.IsClaimed(&bind.CallOpts{}, depositCount, 0)
		require.NoError(t, err)
		require.True(t, isClaimed)
		fmt.Println("birge completed!")
	}
}
