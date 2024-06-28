package service

import (
	"clean-laundry/model"
	"clean-laundry/model/dto"
	"clean-laundry/repository"
	"fmt"
)

// interface - struct - constructor

type ProductService interface {
	FindById(id string) (model.Product, error)
	FindAll(page int, size int) ([]model.Product, dto.Paging, error)
}

type productService struct {
	repo repository.ProductRepository
}

// FindAll implements ProductService.
func (c *productService) FindAll(page int, size int) ([]model.Product, dto.Paging, error) {
	return c.repo.GetAll(page, size)
}

// FindById implements ProductService.
func (c *productService) FindById(id string) (model.Product, error) {
	product, err := c.repo.GetById(id)
	if err != nil {
		return model.Product{}, fmt.Errorf("product with id %s not found", id)
	}

	return product, nil
}

func NewProductService(repositori repository.ProductRepository) ProductService {
	return &productService{repo: repositori}
}
