package handlers

import (
	"database/sql"
	"enigma-laundry/config"
	"enigma-laundry/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

type Response struct {
	Message string          `json:"message"`
	Data    models.Customer `json:"data"`
}

type ResponseList struct {
	Message string            `json:"message"`
	Data    []models.Customer `json:"data"`
}

func CreateCustomer(c *gin.Context) {
	var newCustomer models.Customer
	err := c.ShouldBind(&newCustomer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi phone_number
	phoneNumber := newCustomer.PhoneNumber
	if _, err := strconv.Atoi(phoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must be numeric"})
		return
	}
	if len(phoneNumber) <= 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must be more than 10 digits"})
		return
	}

	//simpan data ke database
	query := "INSERT INTO customers(name, phone_number, address) VALUES ($1,$2,$3) RETURNING id"

	err = db.QueryRow(query, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address).Scan(&newCustomer.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new customer"})
		return
	}

	response := Response{
		Message: "Customer created successfully",
		Data:    newCustomer,
	}

	c.JSON(http.StatusOK, response)
}

func GetAllCustomers(c *gin.Context) {
	searchId := c.Param("id")

	query := "SELECT id, name, phone_number, address FROM customers"

	var rows *sql.Rows
	var err error

	if searchId != "" {
		query += " WHERE id=$1"
		rows, err = db.Query(query, searchId)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var matchedCustomer []models.Customer

	for rows.Next() {
		var customer models.Customer

		err := rows.Scan(&customer.ID, &customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		matchedCustomer = append(matchedCustomer, customer)
	}

	response := ResponseList{
		Message: "Customer updated successfully",
		Data:    matchedCustomer,
	}

	if len(matchedCustomer) > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
	}

}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer id"})
	}

	var updatedCustomer models.Customer
	if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//mengambil data pelanggan yang ada di database berdasarkan ID
	var existingCustomer models.Customer
	query := `SELECT id, name, phone_number, address FROM customers WHERE id=$1;`
	err = db.QueryRow(query, customerId).Scan(&existingCustomer.ID, &existingCustomer.Name, &existingCustomer.PhoneNumber, &existingCustomer.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer"})
		}
		return
	}

	// Pengkondisian: jika setiap field tidak kosong, kita akan update datanya. Jika kosong, kita akan menggunakan data sebelumnya atau partial update
	if strings.TrimSpace(updatedCustomer.Name) != "" {
		existingCustomer.Name = updatedCustomer.Name
	}
	if strings.TrimSpace(updatedCustomer.PhoneNumber) != "" {
		phoneNumber := updatedCustomer.PhoneNumber
		if _, err := strconv.Atoi(phoneNumber); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must be numeric"})
			return
		}
		if len(phoneNumber) <= 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number must be more than 10 digits"})
			return
		}
		existingCustomer.PhoneNumber = updatedCustomer.PhoneNumber
	}
	if strings.TrimSpace(updatedCustomer.Address) != "" {
		existingCustomer.Address = updatedCustomer.Address
	}

	// Update data pelanggan di database
	updateQuery := `UPDATE customers SET name=$1, phone_number=$2, address=$3 WHERE id=$4`
	_, err = db.Exec(updateQuery, existingCustomer.Name, existingCustomer.PhoneNumber, existingCustomer.Address, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "data": existingCustomer})

	// Membuat respons dengan struktur yang diinginkan
	response := Response{
		Message: "Customer updated successfully",
		Data:    existingCustomer,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteCustomerById(c *gin.Context) {
	id := c.Param("id")

	customerId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer Id"})
		return
	}

	query := `DELETE FROM customers WHERE id=$1;`
	_, err = db.Exec(query, customerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete customer"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Customer deleted successfully",
		"data":    "OK",
	})

}
