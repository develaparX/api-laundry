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

	employeesGroup := router.Group("/employees")
	{
		employeesGroup.POST("/", handlers.CreateEmployee)
		employeesGroup.GET("/", handlers.GetAllEmployees)
		employeesGroup.GET("/:id", handlers.GetAllEmployees)
		employeesGroup.PUT("/:id", handlers.UpdateEmployee)
		employeesGroup.DELETE("/:id", handlers.DeleteEmployeeById)
	}

	productsGroup := router.Group("/products")
	{
		productsGroup.POST("/", handlers.CreateProduct)
		productsGroup.GET("/", handlers.GetAllProducts)
		productsGroup.GET("/:id", handlers.GetAllProducts)
		productsGroup.PUT("/:id", handlers.UpdateProduct)
		productsGroup.DELETE("/:id", handlers.DeleteProductById)
	}

	transactionsGroup := router.Group("/transactions")
	{
		transactionsGroup.POST("/", handlers.CreateTransaction)
		transactionsGroup.GET("/:id", handlers.GetTransactionByID)
		transactionsGroup.GET("/", handlers.GetTransactions)
	}

	router.Run(":8080")
}
