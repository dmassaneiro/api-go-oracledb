package models

type Product struct {
	ID    int64   `db:"ID,pk,seq=SEQ_PRODUCTS"`
	Name  string  `db:"NAME"`
	Price float64 `db:"PRICE"`
}

func (Product) TableName() string {
	return "PRODUCTS"
}
