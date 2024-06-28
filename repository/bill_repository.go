package repository

import (
	"clean-laundry/model"
	"database/sql"
	"time"
)

type BillRepository interface {
	Create(payload model.Bill) (model.Bill, error)
}

type billRepository struct {
	db *sql.DB
}

// Create implements BillRepository.
func (b *billRepository) Create(payload model.Bill) (model.Bill, error) {
	transaction, err := b.db.Begin() //transaksi dimulai
	if err != nil {
		return model.Bill{}, err
	}

	//var untuk menampung hasil insertan
	var bill model.Bill
	err = transaction.QueryRow(`INSERT INTO bills(bill_date,customer_id, user_id, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id, bill_date`, time.Now(), payload.Customer.Id, payload.User.Id, time.Now(), time.Now()).Scan(&bill.Id, &bill.BillDate)

	if err != nil { //jika error, di rollback
		return model.Bill{}, transaction.Rollback()
	}

	var billDetails []model.BillDetail

	// perulangan untuk menyimpan kedalam billDetails
	// menggunakan looping karena type datanya berupa slice, jadi satu persatu
	for _, bd := range payload.BillDetails {
		//variable untuk menyimpan
		var billDetail model.BillDetail
		err = transaction.QueryRow(`INSERT INTO bill_details(bill_id, product_id, qty, price, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, qty, price`, bill.Id, bd.Product.Id, bd.Qty, bd.Price, time.Now(), time.Now()).Scan(&billDetail.Id, &billDetail.Qty, &billDetail.Price)
		if err != nil {
			return model.Bill{}, transaction.Rollback()
		}
		billDetail.Product = bd.Product
		billDetails = append(billDetails, billDetail)
	}
	bill.Customer = payload.Customer
	bill.User = payload.User
	bill.BillDetails = billDetails
	if err = transaction.Commit(); err != nil {
		return model.Bill{}, err
	}
	return bill, nil
}

func NewBillRepository(database *sql.DB) BillRepository {
	return &billRepository{db: database}
}
