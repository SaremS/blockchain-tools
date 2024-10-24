package dex

import (
	"fmt"
)

type PoolPriceProcessor interface {
	ProcessTokenPairAndPrice(inputChannel chan TokenPairAndPrice, outputChannel chan<- TokenPairAndPrice)
}

type ConsolePoolPriceProcessor struct{}

func (c *ConsolePoolPriceProcessor) ProcessTokenPairAndPrice(inputChannel chan TokenPairAndPrice, outputChannel chan<- TokenPairAndPrice) {
	for {
		tokenPairAndPrice := <-inputChannel
		token0Symbol := tokenPairAndPrice.tokenPair.token0.GetSymbol()
		token1Symbol := tokenPairAndPrice.tokenPair.token1.GetSymbol()
		price := tokenPairAndPrice.price.GetAsPrice()

		fmt.Printf("%s-%s: %.9f\n", token0Symbol, token1Symbol, price)
		outputChannel <- tokenPairAndPrice
	}
}
