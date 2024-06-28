package service

import (
	"clean-laundry/model"
	"clean-laundry/repository"
	"fmt"
)

// interface - struct - constructor

type CustomerService interface {
	FindById(id string) (model.Customer, error)
	FindAll(page int, size int) ([]model.Customer, error)
}

type customerService struct {
	repo repository.CustomerRepository
}

// FindAll implements CustomerService.
func (c *customerService) FindAll(page int, size int) ([]model.Customer, error) {
	panic("unimplemented")
}

// FindById implements CustomerService.
func (c *customerService) FindById(id string) (model.Customer, error) {
	customer, err := c.repo.GetById(id)
	if err != nil {
		return model.Customer{}, fmt.Errorf("customer with id %s not found", id)
	}

	return customer, nil
}

func NewCustomerService(repositori repository.CustomerRepository) CustomerService {
	return &customerService{repo: repositori}
}
