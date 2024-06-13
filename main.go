package main

import (
	"enigma-laundry/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Tulis kode kamu disini
	router := gin.Default()

	customers := router.Group("/customers")
	{
		customers.POST("/", handlers.CreateCustomer)
	}

	router.Run(":8000")
}
