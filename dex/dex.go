package dex

import (
	"github.com/sarems/blockchain-tools/token"
)

type Dex interface {
	UpdateSwapPrices()
	AddV3Pool(poolAddress string)
	DropPool(poolAddress string)
	GetSwapPriceByPairNames(fromToken, toToken string) float64
	GetSwapPriceByPoolAddress(poolAddress string) float64
	GetTokens() []token.Token
}
