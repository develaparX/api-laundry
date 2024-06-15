package handlers

import (
	"enigma-laundry/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	var request models.CreateTransactionRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Konversi employeeId dan customerId dari string ke int
	employeeID, err := strconv.Atoi(request.EmployeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employeeId harus berupa angka"})
		return
	}

	customerID, err := strconv.Atoi(request.CustomerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customerId harus berupa angka"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	var transactionID int
	err = tx.QueryRow(`INSERT INTO transactions (bill_date, entry_date, finish_date, employee_id, customer_id) 
                       VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		request.BillDate, request.EntryDate, request.FinishDate, employeeID, customerID).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memasukkan data ke transactions"})
		return
	}

	stmt, err := tx.Prepare(`INSERT INTO transaction_details (transaction_id, product_id, product_price, qty) 
                             VALUES ($1, $2, (SELECT price FROM products WHERE id=$2), $3)`)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan statement"})
		return
	}
	defer stmt.Close()

	var billDetails []models.BillDetailResponse

	for _, detail := range request.BillDetails {
		// Konversi productId dari string ke int
		productID, err := strconv.Atoi(detail.ProductID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "productId harus berupa angka"})
			return
		}

		var detailID int
		var productPrice int

		err = stmt.QueryRow(transactionID, productID, detail.Qty).Scan(&detailID, &productPrice)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memasukkan data ke transaction_details"})
			return
		}

		billDetails = append(billDetails, models.BillDetailResponse{
			ID:            detailID,
			TransactionID: transactionID,
			ProductID:     productID,
			ProductPrice:  productPrice,
			Qty:           detail.Qty,
		})
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengkomit transaksi"})
		return
	}

	c.JSON(http.StatusCreated, models.TransactionResponse{
		ID:          transactionID,
		BillDate:    request.BillDate,
		EntryDate:   request.EntryDate,
		FinishDate:  request.FinishDate,
		EmployeeID:  employeeID,
		CustomerID:  customerID,
		BillDetails: billDetails,
	})
}

// func GetTransaction(c *gin.Context) {
// 	transactionID := c.Param("id")

// 	var transaction models.TransactionResponse
// 	err := db.QueryRow(`SELECT id, bill_date, entry_date, finish_date, employee_id, customer_id
//                        FROM transactions
//                        WHERE id = $1`, transactionID).Scan(
// 		&transaction.ID, &transaction.BillDate, &transaction.EntryDate, &transaction.FinishDate,
// 		&transaction.EmployeeID, &transaction.CustomerID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
// 		return
// 	}

// 	rows, err := db.Query(`SELECT td.id, td.transaction_id, td.product_price, td.qty,
//                                   p.id, p.name, p.price, p.unit
//                            FROM transaction_details td
//                            JOIN products p ON td.product_id = p.id
//                            WHERE td.transaction_id = $1`, transactionID)
// 	if err != nil {
// 		log.Fatal("Failed to query transaction details:", err)
// 	}
// 	defer rows.Close()

// 	var billDetails []models.BillDetailResponse
// 	for rows.Next() {
// 		var detail models.BillDetailResponse
// 		var product models.Product

// 		err := rows.Scan(&detail.ID, &detail.TransactionID, &detail.ProductPrice, &detail.Qty,
// 			&product.ID, &product.Name, &product.Price, &product.Unit)
// 		if err != nil {
// 			log.Fatal("Failed to scan transaction detail:", err)
// 		}

// 		detail.Product = product
// 		billDetails = append(billDetails, detail)
// 	}

// 	transaction.BillDetails = billDetails

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Transaction retrieved successfully",
// 		"data":    transaction,
// 	})
// }
