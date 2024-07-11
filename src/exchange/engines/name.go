package engines

type Name int

const (
	UNDEFINED Name = iota
	CURRENCYAPI
	EXCHANGERATEAPI
	FRANKFURTER
)

func (n Name) String() string {
	switch n {
	case CURRENCYAPI:
		return "CurrencyAPI"
	case EXCHANGERATEAPI:
		return "ExchangeRateAPI"
	case FRANKFURTER:
		return "Frankfurter"
	default:
		return "Undefined"
	}
}
