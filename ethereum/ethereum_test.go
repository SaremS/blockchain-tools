package ethereum

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const uniswapV3TokenABI = `[{"constant":true,"inputs":[],"name":"token0","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

func TestGetChainData(t *testing.T) {
	client, err := ethclient.Dial("https://polygon-rpc.com")
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	targetAddress := "0xa374094527e1673a86de625aa59517c5de346d32"
	targetMethod := "token0"

	address := new(common.Address)
	err = GetChainData(client, targetAddress, uniswapV3TokenABI, targetMethod, address)
	if err != nil {
		t.Fatalf("Failed to get token0 address: %v", err)
	}

	expectedAddress := "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270"

	if strings.ToLower(address.Hex()) != expectedAddress {
		t.Fatalf("Expected address to be %s, got %s", expectedAddress, address.Hex())
	}
}
