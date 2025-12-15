package services

import (
	"database/sql"
	"product-api/models"
)

type ProductService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) Create(p *models.Product) error {
	_, err := s.db.Exec(`
		INSERT INTO PRODUCTS (ID, NAME, PRICE)
		VALUES (SEQ_PRODUCTS.NEXTVAL, :1, :2)
	`, p.Name, p.Price)

	return err
}

func (s *ProductService) List() ([]models.Product, error) {
	rows, err := s.db.Query(`
		SELECT ID, NAME, PRICE
		FROM PRODUCTS
		ORDER BY ID
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
