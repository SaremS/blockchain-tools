package token

type Token interface {
	GetTokenName() string
	GetAddressAsString() string
	GetSymbol() string
	GetDecimals() uint8
}
