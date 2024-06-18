package handlers

import (
	"enigma-laundry/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type ResponseTransaction struct {
	Message string               `json:"message"`
	Data    []models.Transaction `json:"data"`
}

type ResponseTransactionByID struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetTransactions(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	productName := c.Query("productName")

	// Parsing tanggal dari format dd-MM-yyyy ke format yyyy-MM-dd
	const layout = "02-01-2006"
	var startDate, endDate string
	var err error

	if startDateStr != "" {
		parsedStartDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Format tanggal startDate tidak valid. Gunakan format dd-MM-yyyy."})
			return
		}
		startDate = parsedStartDate.Format("2006-01-02")
	}

	if endDateStr != "" {
		parsedEndDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Format tanggal endDate tidak valid. Gunakan format dd-MM-yyyy."})
			return
		}
		endDate = parsedEndDate.Format("2006-01-02")
	}

	baseQuery := `
    SELECT
      t.id,
      t.bill_date,
      t.entry_date,
      t.finish_date,
      e.id AS employee_id,
      e.name AS employee_name,
      e.phone_number AS employee_phone_number,
      e.address AS employee_address,
      c.id AS customer_id,
      c.name AS customer_name,
      c.phone_number AS customer_phone_number,
      c.address AS customer_address,
      td.id AS bill_detail_id,
      td.transaction_id AS bill_detail_bill_id,
      p.id AS product_id,
      p.name AS product_name,
      p.price AS product_price,
      p.unit AS product_unit,
      td.product_price,
      td.qty
    FROM transactions t
    JOIN employees e ON t.employee_id = e.id
    JOIN customers c ON t.customer_id = c.id
    JOIN transaction_details td ON t.id = td.transaction_id
    JOIN products p ON td.product_id = p.id
  `

	var conditions []string
	var params []interface{}
	paramIndex := 1

	if startDate != "" {
		conditions = append(conditions, "t.bill_date >= $"+strconv.Itoa(paramIndex))
		params = append(params, startDate)
		paramIndex++
	}

	if endDate != "" {
		conditions = append(conditions, "t.bill_date <= $"+strconv.Itoa(paramIndex))
		params = append(params, endDate)
		paramIndex++
	}

	if productName != "" {
		conditions = append(conditions, "p.name ILIKE '%' || $"+strconv.Itoa(paramIndex)+" || '%'")
		params = append(params, productName)
		paramIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Mulai transaksi
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memulai transaksi"})
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Terjadi kesalahan, transaksi dibatalkan"})
		}
	}()

	rows, err := tx.Query(baseQuery, params...)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	transactions := make(map[int]*models.Transaction)
	for rows.Next() {
		var transaction models.Transaction
		var billDetail models.BillDetail
		var employee models.Employee
		var customer models.Customer
		var product models.Product

		err := rows.Scan(
			&transaction.ID,
			&transaction.BillDate,
			&transaction.EntryDate,
			&transaction.FinishDate,
			&employee.ID,
			&employee.Name,
			&employee.PhoneNumber,
			&employee.Address,
			&customer.ID,
			&customer.Name,
			&customer.PhoneNumber,
			&customer.Address,
			&billDetail.ID,
			&billDetail.BillID,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Unit,
			&billDetail.ProductPrice,
			&billDetail.Qty,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		billDetail.Product = product

		if t, ok := transactions[transaction.ID]; ok {
			t.BillDetails = append(t.BillDetails, billDetail)
			t.TotalBill += billDetail.ProductPrice * billDetail.Qty
		} else {
			transaction.Employee = employee
			transaction.Customer = customer
			transaction.BillDetails = []models.BillDetail{billDetail}
			transaction.TotalBill = billDetail.ProductPrice * billDetail.Qty
			transactions[transaction.ID] = &transaction
		}
	}

	var result []models.Transaction
	for _, t := range transactions {
		result = append(result, *t)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengkonfirmasi transaksi"})
		return
	}

	response := ResponseTransaction{
		Message: "Success",
		Data:    result,
	}

	c.JSON(http.StatusOK, response)
}

func CreateTransaction(c *gin.Context) {
	var req struct {
		BillDate    string `json:"billDate"`
		EntryDate   string `json:"entryDate"`
		FinishDate  string `json:"finishDate"`
		EmployeeID  int    `json:"employeeId"`
		CustomerID  int    `json:"customerId"`
		BillDetails []struct {
			ProductID int `json:"productId"`
			Qty       int `json:"qty"`
		} `json:"billDetails"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to start transaction"})
		return
	}

	var transactionID int
	err = tx.QueryRow(`
		INSERT INTO transactions (bill_date, entry_date, finish_date, employee_id, customer_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, req.BillDate, req.EntryDate, req.FinishDate, req.EmployeeID, req.CustomerID).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transaction"})
		return
	}

	for _, detail := range req.BillDetails {
		var productPrice int
		err := db.QueryRow(`SELECT price FROM products WHERE id = $1`, detail.ProductID).Scan(&productPrice)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch product price"})
			return
		}

		_, err = tx.Exec(`
			INSERT INTO transaction_details (transaction_id, product_id, product_price, qty)
			VALUES ($1, $2, $3, $4)
		`, transactionID, detail.ProductID, productPrice, detail.Qty)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transaction details"})
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created", "data": req})
}

func GetTransactionByID(c *gin.Context) {
	idBill := c.Param("id_bill")

	// Validasi input
	if idBill == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID bill tidak boleh kosong"})
		return
	}

	// Coba konversi ke integer untuk validasi lebih lanjut
	_, err := strconv.Atoi(idBill)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID bill harus berupa angka"})
		return
	}

	// Mulai transaksi
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memulai transaksi"})
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Terjadi kesalahan, transaksi dibatalkan"})
		}
	}()

	query := `
    SELECT
      t.id,
      t.bill_date,
      t.entry_date,
      t.finish_date,
      e.id AS employee_id,
      e.name AS employee_name,
      e.phone_number AS employee_phone_number,
      e.address AS employee_address,
      c.id AS customer_id,
      c.name AS customer_name,
      c.phone_number AS customer_phone_number,
      c.address AS customer_address,
      td.id AS bill_detail_id,
      td.transaction_id AS bill_detail_bill_id,
      p.id AS product_id,
      p.name AS product_name,
      p.price AS product_price,
      p.unit AS product_unit,
      td.product_price,
      td.qty
    FROM transactions t
    JOIN employees e ON t.employee_id = e.id
    JOIN customers c ON t.customer_id = c.id
    JOIN transaction_details td ON t.id = td.transaction_id
    JOIN products p ON td.product_id = p.id
    WHERE t.id = $1
  `

	rows, err := tx.Query(query, idBill)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	var transaction models.Transaction
	transactions := make(map[int]*models.Transaction)
	for rows.Next() {
		var billDetail models.BillDetail
		var employee models.Employee
		var customer models.Customer
		var product models.Product

		err := rows.Scan(
			&transaction.ID,
			&transaction.BillDate,
			&transaction.EntryDate,
			&transaction.FinishDate,
			&employee.ID,
			&employee.Name,
			&employee.PhoneNumber,
			&employee.Address,
			&customer.ID,
			&customer.Name,
			&customer.PhoneNumber,
			&customer.Address,
			&billDetail.ID,
			&billDetail.BillID,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Unit,
			&billDetail.ProductPrice,
			&billDetail.Qty,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		billDetail.Product = product

		if t, ok := transactions[transaction.ID]; ok {
			t.BillDetails = append(t.BillDetails, billDetail)
			t.TotalBill += billDetail.ProductPrice * billDetail.Qty
		} else {
			transaction.Employee = employee
			transaction.Customer = customer
			transaction.BillDetails = []models.BillDetail{billDetail}
			transaction.TotalBill = billDetail.ProductPrice * billDetail.Qty
			transactions[transaction.ID] = &transaction
		}
	}

	if len(transactions) == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaksi tidak ditemukan"})
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengkonfirmasi transaksi"})
		return
	}

	response := ResponseTransactionByID{
		Message: "Success",
		Data:    transactions[transaction.ID],
	}

	c.JSON(http.StatusOK, response)
}
