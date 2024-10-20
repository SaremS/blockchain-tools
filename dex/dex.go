package dex

type Dex interface {
	UpdateSwapPrices()
	AddPool(poolAddress string)
	DropPool(poolAddress string)
	GetSwapPriceByPairNames(fromToken, toToken string) float64
	GetSwapPriceByPoolAddress(poolAddress string) float64
}
