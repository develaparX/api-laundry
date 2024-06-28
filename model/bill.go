package model

import "time"

type Bill struct {
	Id          string       `json:"id"`
	BillDate    string       `json:"billDate"`
	Customer    Customer     `json:"customer"`
	User        User         `json:"user"`
	BillDetails []BillDetail `json:"billDetails"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type BillDetail struct {
	Id        string    `json:"id"`
	BillId    string    `json:"billId"`
	Product   Product   `json:"product"`
	Qty       int       `json:"qty"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
