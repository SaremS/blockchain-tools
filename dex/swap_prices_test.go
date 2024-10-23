package dex

import (
	"testing"
)

func TestQueryTokensAndPrice(t *testing.T) {
	rpcAddress := "https://polygon-rpc.com"
	poolAddress := "0xa374094527e1673a86de625aa59517c5de346d32"

	tokenPairAndPriceChannel := make(chan TokenPairAndPrice)

	go queryTokensAndPrice(rpcAddress, poolAddress, tokenPairAndPriceChannel)

	tokenPairAndPrice := <-tokenPairAndPriceChannel
	token0Name := tokenPairAndPrice.tokenPair.token0.GetTokenName()

	if token0Name != "Wrapped Polygon Ecosystem Token" {
		t.Fatalf("Expected token0 name to be Wrapped Polygon Ecosystem Token, got %s", token0Name)
	}
}

func TestNewSwapPrices(t *testing.T) {
	rpcAddress := "https://polygon-rpc.com"
	poolAddresses := []string{"0xa374094527e1673a86de625aa59517c5de346d32", "0x45dda9cb7c25131df268515131f647d726f50608"}

	poolPriceProcessor := &ConsolePoolPriceProcessor{}

	swapPrices := NewSwapPrices(rpcAddress, poolAddresses, poolPriceProcessor)

	if len(swapPrices.tokensByAddress) != 3 {
		t.Fatalf("Expected 2 tokens by address, got %d", len(swapPrices.tokensByAddress))
	}
}
