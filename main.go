package main

import (
	"fmt"
	"time"

	"github.com/sarems/blockchain-tools/dex"
)

func main() {
	rpcAddress := "https://polygon-rpc.com"
	poolAddresses := []string{"0xa374094527e1673a86de625aa59517c5de346d32", "0x45dda9cb7c25131df268515131f647d726f50608"}

	poolPriceProcessor := &dex.ConsolePoolPriceProcessor{}

	currentTime := time.Now()
	fmt.Println("\nSwap Prices at ", currentTime, ":")
	swapPrices := dex.NewSwapPrices(rpcAddress, poolAddresses, poolPriceProcessor)

	for {
		time.Sleep(5 * time.Second)
		currentTime := time.Now()
		fmt.Println("\nSwap Prices at ", currentTime, ":")

		swapPrices.UpdateAllPrices()
	}
}
