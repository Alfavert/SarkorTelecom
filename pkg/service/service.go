package service

import (
	"SarkorTelekom/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Product interface {
	CreateProduct(products repository.Products) (int, error)
	GetById(productId int) (repository.Products, error)
	GetAll() ([]repository.Products, error)
	Update(productId int, input repository.UpdateProducts) error
	Delete(productId int) error
}

type Service struct {
	Product
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Product: NewProductService(repos.Product),
	}
}
