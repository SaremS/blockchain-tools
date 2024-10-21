package dex

import (
	_ "embed"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/sarems/blockchain-tools/ethereum"
	"github.com/sarems/blockchain-tools/token"
)

//go:embed abi/uniswap_v3_pool.abi.json
var uniswapV3PoolAbi string

type TokenPair struct {
	token0 *token.ERC20Token
	token1 *token.ERC20Token
}

type TokenPairAndPrice struct {
	tokenPair    TokenPair
	sqrtPriceX96 *big.Int
}

type SwapPrices struct {
	tokensByAddress    map[string]*token.Token
	priceByTokenPair   map[TokenPair]*big.Int
	priceByPoolAddress map[string]*big.Int
}

type slot0Data struct {
	SqrtPriceX96               *big.Int
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

//func NewSwapPrice(rpcAddress string, poolAddresses []string) *SwapPrices {
//}

func queryTokensAndPrice(rpcAddress string, poolAddress string, tokenPairAndPriceChannel chan<- TokenPairAndPrice) {
	client, err := ethclient.Dial(rpcAddress)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	token0Address := new(common.Address)
	err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token0", token0Address)
	if err != nil {
		log.Fatalf("Failed to get token0 address: %v", err)
	}

	token0, err := token.NewERC20FromAddressString(client, token0Address.Hex())
	if err != nil {
		log.Fatalf("Failed to create token0: %v", err)
	}

	token1Address := new(common.Address)
	err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token1", token1Address)
	if err != nil {
		log.Fatalf("Failed to get token0 address: %v", err)
	}

	token1, err := token.NewERC20FromAddressString(client, token1Address.Hex())
	if err != nil {
		log.Fatalf("Failed to create token0: %v", err)
	}

	tokenPair := TokenPair{
		token0: token0,
		token1: token1,
	}

	slot0 := new(slot0Data)
	err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "slot0", slot0)
	if err != nil {
		log.Fatalf("Failed to get pool slot0: %v", err)
	}

	tokenPairAndPriceChannel <- TokenPairAndPrice{
		tokenPair:    tokenPair,
		sqrtPriceX96: (*slot0).SqrtPriceX96,
	}
}
