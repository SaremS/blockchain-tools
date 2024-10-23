package token

import (
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestERC20Token(t *testing.T) {
	client, err := ethclient.Dial("https://polygon-rpc.com")
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	targetAddress := "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"

	token, err := NewERC20FromAddressString(client, targetAddress)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	expectedName := "Wrapped Polygon Ecosystem Token"
	if token.GetTokenName() != expectedName {
		t.Fatalf("Expected token name to be %s, got %s", expectedName, token.GetTokenName())
	}

	expectedSymbol := "WPOL"
	if token.GetSymbol() != expectedSymbol {
		t.Fatalf("Expected token symbol to be %s, got %s", expectedSymbol, token.GetSymbol())
	}
}
