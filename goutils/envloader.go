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
