package currencyapi

// Rates field is named the same as base currency.
type response struct {
	Rates map[string]float64 `json:"eur"`
}
