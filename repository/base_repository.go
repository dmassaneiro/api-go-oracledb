
package repository

import "product-api/crud"

type BaseRepository struct {
	crud *crud.Crud
}

func NewBaseRepository(crud *crud.Crud) *BaseRepository {
	return &BaseRepository{crud: crud}
}
