package main

import (
	"enigma-laundry/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tulis kode kamu disini
	router := gin.Default()

	customersGroup := router.Group("/customers")
	{
		customersGroup.POST("/", handlers.CreateCustomer)
		customersGroup.GET("/", handlers.GetAllCustomers)
	}

	router.Run(":8000")
}
