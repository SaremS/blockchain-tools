package dex

import (
	"github.com/sarems/blockchain-tools/token"
)

type UniswapDex struct {
	rpcAddress string
	tokens     []*token.Token
}

func NewUniswapDex(rpcAddress string) *UniswapDex {
	return &UniswapDex{
		rpcAddress: rpcAddress,
	}
}
