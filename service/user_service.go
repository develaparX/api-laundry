package service

import (
	"clean-laundry/model"
	"clean-laundry/model/dto"
	"clean-laundry/repository"
	"clean-laundry/utils"
	"errors"
	"fmt"
)

// interface - struct - constructor

type UserService interface {
	FindById(id string) (model.User, error)
	FindAll(page int, size int) ([]model.User, error)
	CreateNew(payload model.User) (model.User, error)
	FindByUsername(username string) (model.User, error)
	Login(payload dto.LoginDto) (dto.LoginResponseDto, error)
}

type userService struct {
	repo       repository.UserRepository
	jwtService JwtService
}

func (c *userService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	user, err := c.repo.GetByUsername(payload.Username)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("username or password invalid")
	}
	err = utils.ComparePasswordHash(user.Password, payload.Password)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("password incorrect")
	}
	user.Password = ""
	token, err := c.jwtService.GenerateToken(user)
	if err != nil {
		return dto.LoginResponseDto{}, fmt.Errorf("failed create token")
	}
	return token, nil
}

// CreateNew implements UserService.
func (c *userService) CreateNew(payload model.User) (model.User, error) {
	//cek apakah rolenya valid
	if !payload.IsValidRole() {
		return model.User{}, errors.New("role is invalid, must be admin or employee")
	}
	passwordHash, error := utils.EncryptPassword(payload.Password)
	if error != nil {
		return model.User{}, error
	}
	payload.Password = passwordHash
	return c.repo.CreateNew(payload)
}

// FindAll implements UserService.
func (c *userService) FindAll(page int, size int) ([]model.User, error) {
	panic("unimplemented")
}

// GetByUsername implements UserService.
func (u *userService) FindByUsername(username string) (model.User, error) {
	user, err := u.repo.GetByUsername(username)
	if err != nil {
		return model.User{}, fmt.Errorf("user with username %s not found", username)
	}
	return user, nil
}

// FindById implements UserService.
func (c *userService) FindById(id string) (model.User, error) {
	user, err := c.repo.GetById(id)
	if err != nil {
		return model.User{}, fmt.Errorf("user with id %s not found", id)
	}

	return user, nil
}

func NewUserService(repositori repository.UserRepository, jS JwtService) UserService {
	return &userService{repo: repositori, jwtService: jS}
}
