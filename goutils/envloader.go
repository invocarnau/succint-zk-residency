package goutils

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func LoadRPCs(l1ChainID, l2ChainID int) (l1RPC, l2RPC *ethclient.Client, err error) {
	err = godotenv.Load("../.env")
	if err != nil {
		return nil, nil, fmt.Errorf("error loading .env file: %w", err)
	}
	l1URL := os.Getenv(fmt.Sprintf("RPC_%d", l1ChainID))
	if l1URL == "" {
		return nil, nil, fmt.Errorf("L1 URL not found for ChainID %d", l1ChainID)
	}
	l2URL := os.Getenv(fmt.Sprintf("RPC_%d", l2ChainID))
	if l2URL == "" {
		return nil, nil, fmt.Errorf("L2 URL not found for ChainID %d", l2ChainID)
	}
	l1RPC, err = ethclient.Dial(l1URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial L1: %w", err)
	}
	l2RPC, err = ethclient.Dial(l2URL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial L2: %w", err)
	}
	return
}

/////////
/// Quite specific functions to PoS for now, will make it generic going forward
/////////

func LoadEthRpc() (ethRpc *ethclient.Client, err error) {
	err = godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	ethURL := os.Getenv("ETH_RPC_URL")
	if ethURL == "" {
		return nil, fmt.Errorf("invalid ETH_RPC_URL provided")
	}
	ethRpc, err = ethclient.Dial(ethURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ETH_RPC_URL: %w", err)
	}
	return ethRpc, nil
}

func LoadHeimdallEndpoint() (endpoint string, err error) {
	err = godotenv.Load("../.env")
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}
	endpoint = os.Getenv("HEIMDALL_REST_ENDPOINT")
	if endpoint == "" {
		return "", fmt.Errorf("invalid HEIMDALL_REST_ENDPOINT provided")
	}
	return endpoint, nil
}
