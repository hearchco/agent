package engines

import (
	"github.com/hearchco/agent/src/exchange/currency"
)

type Exchanger interface {
	Exchange(base currency.Currency) (currency.Currencies, error)
}
