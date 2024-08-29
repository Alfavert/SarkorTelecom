package repository

import (
	"github.com/jmoiron/sqlx"
)

type Product interface {
	CreateProduct(products Products) (int, error)
	GetById(productId int) (Products, error)
	GetAll() ([]Products, error)
	Update(productId int, input UpdateProducts) error
	Delete(productId int) error
}

type Repository struct {
	Product
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Product: NewProductPostgres(db),
	}
}
