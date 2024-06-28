package repository

import (
	"clean-laundry/model"
	"clean-laundry/model/dto"
	"database/sql"
	"math"
)

// bikin interface
// struct
// constructor

// interface : supaya bisa mudah diakses dari luar
type ProductRepository interface {
	GetById(id string) (model.Product, error)
	GetAll(page int, size int) ([]model.Product, dto.Paging, error)
}

// struct : menaruh depedency/fungsi/library yang akan digunakan
type productRepository struct {
	db *sql.DB
}

// Method GetAll implements ProductRepository.
func (p *productRepository) GetAll(page int, size int) ([]model.Product, dto.Paging, error) {
	//var untu menampung list data
	var listData []model.Product

	//rumus untuk pagination
	skip := (page - 1) * size //misal page 3 size 5 maka output 5 dimulai dari 11 karena skip := (3-1)*5 maka skip offset 10 dan output 11-15

	rows, err := p.db.Query(`SELECT * FROM products LIMIT $1 OFFSET $2`, size, skip)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	totalRows := 0
	err = p.db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&totalRows)
	if err != nil {
		return nil, dto.Paging{}, err
	}

	for rows.Next() {
		var product model.Product

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Type,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		listData = append(listData, product)
	}

	paging := dto.Paging{
		Page:       page,
		Size:       size,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64((totalRows) / (size)))),
	}
	return listData, paging, nil

}

// Method GetById implements ProductRepository.
func (p *productRepository) GetById(id string) (model.Product, error) {
	var product model.Product

	err := p.db.QueryRow(`SELECT id,name, price, type, created_at FROM products WHERE id=$1`, id).Scan(&product.Id, &product.Name, &product.Price, &product.Type, &product.CreatedAt) // .scan chaining method,
	if err != nil {
		return model.Product{}, err
	}
	return product, nil
}

// constructur : function yang diakses pertama kali ketika dijalankan
// memasukkan return di ProductRepository ke &productRepository
// istilah mudahnya untuk memaskukkan method didalam interface ProductRepository kedalam database(struct productRepository)
func NewProductRepository(database *sql.DB) ProductRepository { // mereturn product repo
	return &productRepository{db: database}

}
