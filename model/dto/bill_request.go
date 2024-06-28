package dto

import "clean-laundry/model"

// DTO (Data transfer object) : sebagai wadah sementara/untuk transfer sajas
// DTO untuk mengirim request bill
type BillRequest struct {
	Id          string             `json:"id"`
	CustomerId  string             `json:"customerId"`
	UserId      string             `json:"userId"`
	BillDetails []model.BillDetail `json:"billDetails"`
}
