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
		token0Name := tokenPairAndPrice.tokenPair.token0.GetTokenName()
		token1Name := tokenPairAndPrice.tokenPair.token1.GetTokenName()
		price := tokenPairAndPrice.price.GetAsPrice()

		fmt.Printf("Token pair: %s and %s - pool price %.6f\n", token0Name, token1Name, price)
		outputChannel <- tokenPairAndPrice
	}
}
