package handlers

import (
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
