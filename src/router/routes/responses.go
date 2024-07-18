package routes

import (
	"github.com/hearchco/agent/src/exchange/currency"
	"github.com/hearchco/agent/src/search/result"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}

type responseBase struct {
	Version  string `json:"version"`
	Duration int64  `json:"duration"`
}

type ResultsResponse struct {
	responseBase

	Results []result.ResultOutput `json:"results"`
}

type SuggestionsResponse struct {
	responseBase

	Suggestions []result.Suggestion `json:"suggestions"`
}

type ExchangeResponse struct {
	responseBase

	Base   currency.Currency `json:"base"`
	From   currency.Currency `json:"from"`
	To     currency.Currency `json:"to"`
	Amount float64           `json:"amount"`
	Result float64           `json:"result"`
}

type CurrenciesResponse struct {
	responseBase

	Base       currency.Currency   `json:"base"`
	Currencies currency.Currencies `json:"currencies"`
}
