package repository

import (
	"clean-laundry/model"
	"database/sql"
)

// bikin interface
// struct
// constructor

// interface : supaya bisa mudah diakses dari luar
type CustomerRepository interface {
	GetById(id string) (model.Customer, error)
	GetAll(page int, size int) ([]model.Customer, error)
}

// struct : menaruh depedency/fungsi/library yang akan digunakan
type customerRepository struct {
	db *sql.DB
}

// Method GetAll implements CustomerRepository.
func (p *customerRepository) GetAll(page int, size int) ([]model.Customer, error) {
	panic("unimplemented")
}

// Method GetById implements CustomerRepository.
func (p *customerRepository) GetById(id string) (model.Customer, error) {
	var customer model.Customer

	err := p.db.QueryRow(`SELECT id, name, phone_number, address, created_at FROM customers WHERE id=$1`, id).Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address, &customer.CreatedAt) // .scan chaining method,
	if err != nil {
		return model.Customer{}, err
	}
	return customer, nil
}

// constructur : function yang diakses pertama kali ketika dijalankan
// memasukkan return di CustomerRepository ke &customerRepository
// istilah mudahnya untuk memaskukkan method didalam interface CustomerRepository kedalam database(struct customerRepository)
func NewCustomerRepository(database *sql.DB) CustomerRepository { // mereturn customer repo
	return &customerRepository{db: database}

}
