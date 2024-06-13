package handlers

import (
	"database/sql"
	"enigma-laundry/config"
	"enigma-laundry/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

func CreateCustomer(c *gin.Context) {
	var newCustomer models.Customer
	err := c.ShouldBind(&newCustomer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//simpan data ke database
	query := "INSERT INTO customers(name, phone_number, address) VALUES ($1,$2,$3) RETURNING id"

	err = db.QueryRow(query, newCustomer.Name, newCustomer.PhoneNumber, newCustomer.Address).Scan(&newCustomer.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    newCustomer,
	})
}

func GetAllCustomers(c *gin.Context) {
	searchId := c.Query("id")

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

	if len(matchedCustomer) > 0 {
		c.JSON(http.StatusOK, matchedCustomer)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Customer not found",
		})
	}

}
