package repository

import "product-api/models"

type ProductRepository struct {
	*BaseRepository
}

func NewProductRepository(base *BaseRepository) *ProductRepository {
	return &ProductRepository{BaseRepository: base}
}

func (r *ProductRepository) Create(p models.Product) (models.Product, error) {
	var id int64

	err := r.crud.CreateStructReturningID(
		p.TableName(),
		p,
		&id,
	)
	if err != nil {
		return p, err
	}

	p.ID = id
	return p, nil
}

func (r *ProductRepository) List() ([]models.Product, error) {
	var products []models.Product

	var p models.Product
	err := r.crud.ListStruct(p.TableName(), &products)

	return products, err
}

func (r *ProductRepository) Update(p models.Product) error {
	return r.crud.UpdateStruct(p.TableName(), p)
}

func (r *ProductRepository) Delete(id int64) error {
	p := models.Product{
		ID: id,
	}
	return r.crud.DeleteByPK(p.TableName(), &p)
}

func (r *ProductRepository) FindByID(id int64) (models.Product, error) {
	var p models.Product
	p.ID = id

	err := r.crud.FindByID(p.TableName(), &p)
	return p, err
}
