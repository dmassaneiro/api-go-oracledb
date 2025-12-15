package facade

import (
	"errors"
	"product-api/models"
	"product-api/repository"
)

type ProductFacade struct {
	repo *repository.ProductRepository
}

func NewProductFacade(repo *repository.ProductRepository) *ProductFacade {
	return &ProductFacade{repo: repo}
}

func (f *ProductFacade) Create(p models.Product) (models.Product, error) {
	if p.Name == "" {
		return p, errors.New("nome é obrigatório")
	}
	if p.Price <= 0 {
		return p, errors.New("preço inválido")
	}

	return f.repo.Create(p)
}

func (f *ProductFacade) List() ([]models.Product, error) {
	return f.repo.List()
}

func (f *ProductFacade) FindByID(id int64) (models.Product, error) {
	return f.repo.FindByID(id)
}

func (f *ProductFacade) Update(id int64, p models.Product) (models.Product, error) {
	if p.Name == "" {
		return p, errors.New("nome é obrigatório")
	}
	if p.Price <= 0 {
		return p, errors.New("preço inválido")
	}

	p.ID = id
	err := f.repo.Update(p)
	return p, err
}

func (f *ProductFacade) Delete(id int64) error {
	return f.repo.Delete(id)
}
