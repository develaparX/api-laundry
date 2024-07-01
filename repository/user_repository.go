package repository

import (
	"clean-laundry/model"
	"database/sql"
	"fmt"
	"time"
)

// bikin interface
// struct
// constructor

// interface : supaya bisa mudah diakses dari luar
type UserRepository interface {
	GetById(id string) (model.User, error)
	GetAll(page int, size int) ([]model.User, error)
	CreateNew(payload model.User) (model.User, error)
	GetByUsername(username string) (model.User, error)
}

// struct : menaruh depedency/fungsi/library yang akan digunakan
type userRepository struct {
	db *sql.DB
}

// CreateNew implements UserRepository.
func (p *userRepository) CreateNew(payload model.User) (model.User, error) {
	var user model.User
	err := p.db.QueryRow(`INSERT INTO users (name, email,username, password, role, updated_at) VALUES($1, $2,$3,$4,$5,$6) RETURNING id, name, email, username, role, created_at`, payload.Name, payload.Email, payload.Username, payload.Password, payload.Role, time.Now()).Scan(&user.Id, &user.Name, &user.Email, &user.Username, &user.Role, &user.CreatedAt)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// Method GetAll implements UserRepository.
func (p *userRepository) GetAll(page int, size int) ([]model.User, error) {
	panic("unimplemented")
}

// GetByUsername implements UserRepository.
func (u *userRepository) GetByUsername(username string) (model.User, error) {
	var user model.User

	query := `
        SELECT id, name, email, username, password, role, created_at, updated_at
        FROM users
        WHERE username = $1
    `

	err := u.db.QueryRow(query, username).Scan(
		&user.Id, &user.Name, &user.Email, &user.Username, &user.Password,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, fmt.Errorf("user with username %s not found", username)
		}
		return model.User{}, err
	}

	return user, nil
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
