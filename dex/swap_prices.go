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
	token0      token.Token
	token1      token.Token
	poolAddress string
}

func (t *TokenPair) SwapTokens() *TokenPair {
	token0 := t.token1
	token1 := t.token0

	return &TokenPair{
		token0:      token0,
		token1:      token1,
		poolAddress: t.poolAddress,
	}
}

type TokenPairAndPrice struct {
	tokenPair TokenPair
	price     Price
}

type SwapPrices struct {
	tokensByAddress     map[string]token.Token
	tokensBySymbol      map[string]token.Token
	priceByTokenPair    map[TokenPair]Price
	tokensByPoolAddress map[string]TokenPair
	poolPriceProcessor  PoolPriceProcessor
	rpcAddress          string
}

func (s *SwapPrices) GetTokenBySymbol(symbol string) token.Token {
	return s.tokensBySymbol[symbol]
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

func NewSwapPrices(rpcAddress string, poolAddresses []string, poolPriceProcessor PoolPriceProcessor) *SwapPrices {
	tokenPairAndPriceOutputChannel := make(chan TokenPairAndPrice)
	for _, poolAddress := range poolAddresses {
		tokenPairAndPriceInputChannel := make(chan TokenPairAndPrice)
		go queryTokensAndPrice(rpcAddress, poolAddress, tokenPairAndPriceInputChannel)
		go poolPriceProcessor.ProcessTokenPairAndPrice(tokenPairAndPriceInputChannel, tokenPairAndPriceOutputChannel)
	}
	tokenPairAndPrices := make([]TokenPairAndPrice, 0)
	for range poolAddresses {
		tokenPairWithPrice := <-tokenPairAndPriceOutputChannel
		tokenPairAndPrices = append(tokenPairAndPrices, tokenPairWithPrice)
	}
	return tokenAndPairPricesToSwapPrices(tokenPairAndPrices, poolPriceProcessor, rpcAddress)
}

func (s *SwapPrices) UpdateAllPrices() {
	tokenPairAndPriceOutputChannel := make(chan TokenPairAndPrice)
	for _, tokenPair := range s.tokensByPoolAddress {
		tokenPairAndPriceInputChannel := make(chan TokenPairAndPrice)
		go queryTokensAndPrice(s.rpcAddress, tokenPair.poolAddress, tokenPairAndPriceInputChannel)
		go s.poolPriceProcessor.ProcessTokenPairAndPrice(tokenPairAndPriceInputChannel, tokenPairAndPriceOutputChannel)
	}
	tokenPairAndPrices := make([]TokenPairAndPrice, 0)
	for range s.tokensByPoolAddress {
		tokenPairWithPrice := <-tokenPairAndPriceOutputChannel
		tokenPairAndPrices = append(tokenPairAndPrices, tokenPairWithPrice)
	}
	s.updatePrices(tokenPairAndPrices)
}

func (s *SwapPrices) updatePrices(tokenPairAndPrices []TokenPairAndPrice) {
	for _, tokenPairAndPrice := range tokenPairAndPrices {
		tokenPair := tokenPairAndPrice.tokenPair
		price := tokenPairAndPrice.price
		s.priceByTokenPair[tokenPair] = price
	}
}

func tokenAndPairPricesToSwapPrices(tokenPairAndPrices []TokenPairAndPrice, poolPriceProcessor PoolPriceProcessor, rpcAddress string) *SwapPrices {
	tokensByAddress := make(map[string]token.Token)
	tokensBySymbol := make(map[string]token.Token)
	priceByTokenPair := make(map[TokenPair]Price)
	tokensByPoolAddress := make(map[string]TokenPair)
	for _, tokenPairAndPrice := range tokenPairAndPrices {
		tokenPair := tokenPairAndPrice.tokenPair
		price := tokenPairAndPrice.price
		tokensByAddress[tokenPair.token0.GetAddressAsString()] = tokenPair.token0
		tokensByAddress[tokenPair.token1.GetAddressAsString()] = tokenPair.token1
		tokensBySymbol[tokenPair.token0.GetSymbol()] = tokenPair.token0
		tokensBySymbol[tokenPair.token1.GetSymbol()] = tokenPair.token1
		priceByTokenPair[tokenPair] = price
		tokensByPoolAddress[tokenPair.poolAddress] = tokenPair
	}
	return &SwapPrices{
		tokensByAddress:     tokensByAddress,
		tokensBySymbol:      tokensBySymbol,
		priceByTokenPair:    priceByTokenPair,
		tokensByPoolAddress: tokensByPoolAddress,
		poolPriceProcessor:  poolPriceProcessor,
		rpcAddress:          rpcAddress,
	}
}

func queryTokensAndPrice(rpcAddress string, poolAddress string, tokenPairAndPriceChannel chan<- TokenPairAndPrice) {
	retriesDone := 0

	for retriesDone < 5 {
		client, err := ethclient.Dial(rpcAddress)
		if err != nil {
			log.Printf("Failed to connect to the Ethereum client: %v", err)
			continue
		}

		token0Address := new(common.Address)
		err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token0", token0Address)
		if err != nil {
			log.Printf("Failed to get token0 address: %v", err)
			continue
		}

		token0, err := token.NewERC20FromAddressString(client, token0Address.Hex())
		if err != nil {
			log.Printf("Failed to create token0: %v", err)
			continue
		}

		token1Address := new(common.Address)
		err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "token1", token1Address)
		if err != nil {
			log.Printf("Failed to get token0 address: %v", err)
			continue
		}

		token1, err := token.NewERC20FromAddressString(client, token1Address.Hex())
		if err != nil {
			log.Printf("Failed to create token0: %v", err)
			continue
		}

		tokenPair := TokenPair{
			token0:      token0,
			token1:      token1,
			poolAddress: poolAddress,
		}

		slot0 := new(slot0Data)
		err = ethereum.GetChainData(client, poolAddress, uniswapV3PoolAbi, "slot0", slot0)
		if err != nil {
			log.Printf("Failed to get pool slot0: %v", err)
			continue
		}

		sqrtPriceX96 := slot0.SqrtPriceX96
		price := NewUniswapV3Price(sqrtPriceX96, token0.GetDecimals(), token1.GetDecimals())

		tokenPairAndPriceChannel <- TokenPairAndPrice{
			tokenPair: tokenPair,
			price:     price,
		}

		return
	}

	log.Fatalf("Failed to query tokens and price for pool %s after 5 retries", poolAddress)
}
