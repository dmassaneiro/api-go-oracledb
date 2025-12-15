package mappers

import (
	req "product-api/dto/request"
	res "product-api/dto/response"
	"product-api/models"
)

func ToProductResponse(p models.Product) res.ProductResponseDTO {
	return res.ProductResponseDTO{
		ID:    p.ID,
		Name:  p.Name,
		Price: p.Price,
	}
}

func ToProductModel(req req.ProductRequestDTO) models.Product {
	return models.Product{
		Name:  req.Name,
		Price: req.Price,
	}
}

func ToProductResponseList(products []models.Product) []res.ProductResponseDTO {
	list := make([]res.ProductResponseDTO, 0, len(products))

	for _, p := range products {
		list = append(list, ToProductResponse(p))
	}

	return list
}
