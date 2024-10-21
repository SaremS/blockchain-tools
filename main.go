package main

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	uniswapV3PoolABI  = ``
	uniswapV3TokenABI = ``
)

const (
	poolAddress = "0x45dda9cb7c25131df268515131f647d726f50608"
)

const ethRPC = "https://polygon-rpc.com"

type slot0Data struct {
	SqrtPriceX96               *big.Int // uint160
	Tick                       *big.Int
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
	Unlocked                   bool
}

func sqrtPriceX96ToPrice(sqrtPriceX96 *big.Int, decimalsToken0, decimalsToken1 uint8) float64 {
	scaleFactor := new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)

	priceX192 := new(big.Int).Mul(sqrtPriceX96, sqrtPriceX96)

	price := new(big.Float).Quo(
		new(big.Float).SetInt(priceX192),
		new(big.Float).SetInt(scaleFactor),
	)

	priceFloat, _ := price.Float64()

	decimalDiff := int(decimalsToken0) - int(decimalsToken1)
	scaler := math.Pow(10.0, float64(decimalDiff))

	return priceFloat * scaler
}

func main() {
	var data slot0Data
	var addressToken0 common.Address
	var addressToken1 common.Address
	var decimalsToken0 uint8
	var decimalsToken1 uint8

	err := GetChainData(ethRPC, poolAddress, uniswapV3PoolABI, "slot0", &data)
	if err != nil {
		fmt.Printf("Error while fetching slot0 data: %v\n", err)
		return
	}

	err = GetChainData(ethRPC, poolAddress, uniswapV3TokenABI, "token0", &addressToken0)
	if err != nil {
		fmt.Printf("Error while fetching slot0 data: %v\n", err)
		return
	}

	err = GetChainData(ethRPC, poolAddress, uniswapV3TokenABI, "token1", &addressToken1)
	if err != nil {
		fmt.Printf("Error while fetching slot0 data: %v\n", err)
		return
	}

	err = GetChainData(ethRPC, addressToken0.Hex(), erc20ABI, "decimals", &decimalsToken0)
	if err != nil {
		fmt.Printf("Error while fetching token decimals: %v\n", err)
		return
	}

	err = GetChainData(ethRPC, addressToken1.Hex(), erc20ABI, "decimals", &decimalsToken1)
	if err != nil {
		fmt.Printf("Error while fetching token decimals: %v\n", err)
		return
	}

	price := sqrtPriceX96ToPrice(data.SqrtPriceX96, decimalsToken0, decimalsToken1)
	fmt.Printf("Price: %f\n", price)
}
