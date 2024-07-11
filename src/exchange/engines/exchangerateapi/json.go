package exchangerateapi

type response struct {
	Rates map[string]float64 `json:"rates"`
}
