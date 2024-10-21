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
	token0Name := (*tokenPairAndPrice.tokenPair.token0).GetTokenName()

	if token0Name != "Wrapped Polygon Ecosystem Token" {
		t.Fatalf("Expected token0 name to be Wrapped Polygon Ecosystem Token, got %s", token0Name)
	}
}
