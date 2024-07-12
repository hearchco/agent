package currencyapi

import (
	"github.com/hearchco/agent/src/exchange/currency"
)

type Exchange struct{}

func New() Exchange {
	return Exchange{}
}

func (e Exchange) apiUrlWithBaseCurrency(base currency.Currency) string {
	return apiUrl + "/" + base.Lower() + ".json"
}
