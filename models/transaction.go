package models

type BillDetail struct {
	ProductID string `json:"productId"`
	Qty       int    `json:"qty"`
}

type CreateTransactionRequest struct {
	BillDate    string       `json:"billDate"`
	EntryDate   string       `json:"entryDate"`
	FinishDate  string       `json:"finishDate"`
	EmployeeID  string       `json:"employeeId"`
	CustomerID  string       `json:"customerId"`
	BillDetails []BillDetail `json:"billDetails"`
}

type BillDetailResponse struct {
	ID            int `json:"id"`
	TransactionID int `json:"transactionId"`
	ProductID     int `json:"productId"`
	ProductPrice  int `json:"productPrice"`
	Qty           int `json:"qty"`
}

type TransactionResponse struct {
	ID          int                  `json:"id"`
	BillDate    string               `json:"billDate"`
	EntryDate   string               `json:"entryDate"`
	FinishDate  string               `json:"finishDate"`
	EmployeeID  int                  `json:"employeeId"`
	CustomerID  int                  `json:"customerId"`
	BillDetails []BillDetailResponse `json:"billDetails"`
}
