package service

import (
	"SarkorTelekom/pkg/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (p *ProductService) CreateProduct(products repository.Products) (int, error) {
	//_, err := p.repo.GetById(products.Id)
	//if err != nil {
	//	// list does not exists or does not belongs to user
	//	return 0, err
	//}

	return p.repo.CreateProduct(products)
}

func (p *ProductService) GetById(productId int) (repository.Products, error) {
	return p.repo.GetById(productId)
}

func (p *ProductService) GetAll() ([]repository.Products, error) {
	return p.repo.GetAll()
}

func (p *ProductService) Update(productId int, input repository.UpdateProducts) error {
	return p.repo.Update(productId, input)
}

func (p *ProductService) Delete(productId int) error {
	return p.repo.Delete(productId)
}
