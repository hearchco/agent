package routes

type ErrorResponse struct {
	Message string `json:"message"`
	Value   string `json:"value"`
}
