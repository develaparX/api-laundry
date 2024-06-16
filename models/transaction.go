package models

type BillDetail struct {
	ID           int     `json:"id"`
	BillID       int     `json:"billId"`
	Product      Product `json:"product"`
	ProductPrice int     `json:"productPrice"`
	Qty          int     `json:"qty"`
}

type Transaction struct {
	ID          int          `json:"id"`
	BillDate    string       `json:"billDate"`
	EntryDate   string       `json:"entryDate"`
	FinishDate  string       `json:"finishDate"`
	Employee    Employee     `json:"employee"`
	Customer    Customer     `json:"customer"`
	BillDetails []BillDetail `json:"billDetails"`
	TotalBill   int          `json:"totalBill"`
}
