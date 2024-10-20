package token

import (
  "embed"

  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/sarems/blockchain-tools/ethereum"
)

//go:embed abi/erc20.abi.json
var ERC20ABI string

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

func NewERC20FromAddressString(client *ethclient.Client, address string) (*ERC20Token, error) {}
  var name *string = new(string)
  err := ethereum.GetChainData(client, address, ERC20ABI, "name", name)
  if err != nil {
    return nil, err
  }

  var decimals *uint = new(uint)
  err = ethereum.GetChainData(client, address, erc20ABI, "decimals", decimals)
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
