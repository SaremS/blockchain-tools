package token

type Token interface {
	GetTokenName() string
	GetAddressAsString() string
	GetDecimals() uint
}
