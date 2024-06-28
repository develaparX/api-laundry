package service

import (
	"clean-laundry/model"
	"clean-laundry/repository"
	"fmt"
)

// interface - struct - constructor

type UserService interface {
	FindById(id string) (model.User, error)
	FindAll(page int, size int) ([]model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// FindAll implements UserService.
func (c *userService) FindAll(page int, size int) ([]model.User, error) {
	panic("unimplemented")
}

// FindById implements UserService.
func (c *userService) FindById(id string) (model.User, error) {
	user, err := c.repo.GetById(id)
	if err != nil {
		return model.User{}, fmt.Errorf("user with id %s not found", id)
	}

	return user, nil
}

func NewUserService(repositori repository.UserRepository) UserService {
	return &userService{repo: repositori}
}
