package dex

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/sarems/blockchain-tools/ethereum"
	"github.com/sarems/blockchain-tools/token"
)

//go:embed abi/uniswap_v3_pool.abi.json
var uniswapV3PoolAbi string

type TokenPair struct {
	token0 *token.Token
	token1 *token.Token
}

type TokenPairAndPrice struct {
	tokenPair TokenPair
	price     *big.Int
}

type SwapPrices struct {
	tokensByAddress    map[string]*token.Token
	priceByTokenPair   map[TokenPair]*big.Int
	priceByPoolAddress map[string]*big.Int
}

func NewSwapPrice(rpcAddress string, poolAddresses []string) *SwapPrices {
}

func (s *SwapPrices) queryTokensAndPrice(rpcAddress string, poolAddress string, tokenPairAndPrice chan<- TokenPairAndPrice) {
	client, err := ethclient.Dial(rpcAddress)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var token0Address *string = new(string)
	err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token0", token0Address)
	if err != nil {
		log.Fatalf("Failed to get token0 address: %v", err)
	}

	token0, err := token.NewERC20FromAddressString(client, *token0Address)
	if err != nil {
		log.Fatalf("Failed to create token0: %v", err)
	}

	var token1Address *string = new(string)
	err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token0", token1Address)
	if err != nil {
		log.Fatalf("Failed to get token0 address: %v", err)
	}

	token1, err := token.NewERC20FromAddressString(client, *token1Address)
	if err != nil {
		log.Fatalf("Failed to create token0: %v", err)
	}

	tokenPair := TokenPair{
		token0: token0,
		token1: token1,
	}

	var price *big.Int = new(big.Int)
	err := ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "price", price)
	if err != nil {
		log.Fatalf("Failed to get token0 address: %v", err)
	}

	return tokenPair, price
}
