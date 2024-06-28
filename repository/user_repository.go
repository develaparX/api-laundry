package repository

import (
	"clean-laundry/model"
	"database/sql"
)

// bikin interface
// struct
// constructor

// interface : supaya bisa mudah diakses dari luar
type UserRepository interface {
	GetById(id string) (model.User, error)
	GetAll(page int, size int) ([]model.User, error)
}

// struct : menaruh depedency/fungsi/library yang akan digunakan
type userRepository struct {
	db *sql.DB
}

// Method GetAll implements UserRepository.
func (p *userRepository) GetAll(page int, size int) ([]model.User, error) {
	panic("unimplemented")
}

// Method GetById implements UserRepository.
func (p *userRepository) GetById(id string) (model.User, error) {
	var user model.User

	err := p.db.QueryRow(`SELECT id, name, email, username, password, role, created_at FROM users WHERE id=$1`, id).Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Password, &user.Role, &user.CreatedAt) // .scan chaining method,
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// constructur : function yang diakses pertama kali ketika dijalankan
// memasukkan return di UserRepository ke &userRepository
// istilah mudahnya untuk memaskukkan method didalam interface UserRepository kedalam database(struct userRepository)
func NewUserRepository(database *sql.DB) UserRepository { // mereturn user repo
	return &userRepository{db: database}

}
