package token

import (
	_ "embed"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sarems/blockchain-tools/ethereum"
)

//go:embed abi/erc20.abi.json
var erc20Abi string

type ERC20Token struct {
	name     string
	address  string
	symbol   string
	decimals uint8
}

func NewERC20Token(name, address, symbol string, decimals uint8) *ERC20Token {
	return &ERC20Token{
		name:     name,
		address:  address,
		symbol:   symbol,
		decimals: decimals,
	}
}

func NewERC20FromAddressString(client *ethclient.Client, address string) (*ERC20Token, error) {
	name := new(string)
	err := ethereum.GetChainData(client, address, erc20Abi, "name", name)
	if err != nil {
		return nil, err
	}

	decimals := new(uint8)
	err = ethereum.GetChainData(client, address, erc20Abi, "decimals", decimals)
	if err != nil {
		return nil, err
	}

	symbol := new(string)
	err = ethereum.GetChainData(client, address, erc20Abi, "symbol", symbol)
	if err != nil {
		return nil, err
	}

	return NewERC20Token(*name, address, *symbol, *decimals), nil
}

func (e *ERC20Token) GetTokenName() string {
	return e.name
}

func (e *ERC20Token) GetAddressAsString() string {
	return e.address
}

func (e *ERC20Token) GetSymbol() string {
	return e.symbol
}

func (e *ERC20Token) GetDecimals() uint8 {
	return e.decimals
}
