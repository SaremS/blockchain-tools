package token

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sarems/blockchain-tools/ethereum"
)

//go:embed abi/erc20.abi.json
var erc20Abi string

type ERC20Token struct {
	name     string
	address  string
	decimals uint
}

func NewERC20Token(name string, address string, decimals uint) *ERC20Token {
	return &ERC20Token{
		name:     name,
		address:  address,
		decimals: decimals,
	}
}

func NewERC20FromAddressString(client *ethclient.Client, address string) (*ERC20Token, error) {
	name := new(string)
	err := ethereum.GetChainData(client, address, erc20Abi, "name", name)
	if err != nil {
		return nil, err
	}

	decimals := new(uint)
	err = ethereum.GetChainData(client, address, erc20Abi, "decimals", decimals)
	if err != nil {
		return nil, err
	}

	return NewERC20Token(*name, address, *decimals), nil
}

func (e *ERC20Token) GetTokenName() string {
	return e.name
}

func (e *ERC20Token) GetAddressAsString() string {
	return e.address
}

func (e *ERC20Token) GetDecimals() uint {
	return e.decimals
}
