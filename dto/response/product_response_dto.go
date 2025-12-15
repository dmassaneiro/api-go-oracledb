package response

type ProductResponseDTO struct {
	ID    int64   `json:"id" example:"1"`
	Name  string  `json:"name" example:"Teclado Mecânico"`
	Price float64 `json:"price" example:"499.90"`
}
