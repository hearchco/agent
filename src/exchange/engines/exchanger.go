package engines

import (
	"github.com/hearchco/agent/src/exchange/currency"
)

type Exchanger interface {
	Exchange(base currency.Currency) (map[currency.Currency]float64, error)
}
