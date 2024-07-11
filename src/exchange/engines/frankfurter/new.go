package frankfurter

import (
	"github.com/hearchco/agent/src/exchange/currency"
)

type Exchange struct {
	apiUrl string
}

func New() Exchange {
	return Exchange{apiUrl}
}

func (e Exchange) apiUrlWithBaseCurrency(base currency.Currency) string {
	return apiUrl + "?from=" + base.String()
}
