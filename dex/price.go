package dex

import (
	"math"
	"math/big"
)

type Price interface {
	GetAsSqrtPriceX96() *big.Int
	GetAsPrice() float64
}

type UniswapV3Price struct {
	sqrtPriceX96   *big.Int
	decimalsToken0 uint8
	decimalsToken1 uint8
}

func NewUniswapV3Price(sqrtPriceX96 *big.Int, decimalsToken0, decimalsToken1 uint8) *UniswapV3Price {
	return &UniswapV3Price{
		sqrtPriceX96:   sqrtPriceX96,
		decimalsToken0: decimalsToken0,
		decimalsToken1: decimalsToken1,
	}
}

func (u *UniswapV3Price) GetAsSqrtPriceX96() *big.Int {
	return u.sqrtPriceX96
}

func (u *UniswapV3Price) GetAsPrice() float64 {
	scaleFactor := new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)

	priceX192 := new(big.Int).Mul(u.sqrtPriceX96, u.sqrtPriceX96)

	price := new(big.Float).Quo(
		new(big.Float).SetInt(priceX192),
		new(big.Float).SetInt(scaleFactor),
	)

	priceFloat, _ := price.Float64()

	decimalDiff := int(u.decimalsToken0) - int(u.decimalsToken1)
	scaler := math.Pow(10.0, float64(decimalDiff))

	return priceFloat * scaler
}
