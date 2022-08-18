package response

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Status  string `json:"status"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}
