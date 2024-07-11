package frankfurter

// Rates doesn't include the base currency.
type response struct {
	Rates map[string]float64 `json:"rates"`
}
