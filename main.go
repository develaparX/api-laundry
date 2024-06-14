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
		customersGroup.GET("/:id", handlers.GetAllCustomers)
		customersGroup.PUT("/:id", handlers.UpdateCustomer)
		customersGroup.DELETE("/:id", handlers.DeleteCustomerById)
	}

	router.Run(":8000")
}
