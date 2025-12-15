package response

type ErrorResponseDTO struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}
