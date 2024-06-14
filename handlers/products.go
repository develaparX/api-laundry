package handlers

import (
	"database/sql"
	"enigma-laundry/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseProduct struct {
	Message string         `json:"message"`
	Data    models.Product `json:"data"`
}

type ResponseListProduct struct {
	Message string           `json:"message"`
	Data    []models.Product `json:"data"`
}

func CreateProduct(c *gin.Context) {
	var newProduct models.Product
	err := c.ShouldBind(&newProduct)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//simpan data ke database
	query := "INSERT INTO Products(name, price, unit) VALUES ($1,$2,$3) RETURNING id"

	err = db.QueryRow(query, newProduct.Name, newProduct.Price, newProduct.Unit).Scan(&newProduct.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new Product"})
		return
	}

	response := ResponseProduct{
		Message: "Product created successfully",
		Data:    newProduct,
	}

	c.JSON(http.StatusCreated, response)
}

func GetAllProducts(c *gin.Context) {
	searchId := c.Param("id")

	query := "SELECT id, name, price, unit FROM products"

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

	var matchedProduct []models.Product

	for rows.Next() {
		var Product models.Product

		err := rows.Scan(&Product.ID, &Product.Name, &Product.Price, &Product.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		matchedProduct = append(matchedProduct, Product)
	}

	response := ResponseListProduct{
		Message: "Product retreived successfully",
		Data:    matchedProduct,
	}

	if len(matchedProduct) > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
	}

}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	ProductId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product id"})
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//mengambil data pelanggan yang ada di database berdasarkan ID
	var existingProduct models.Product
	query := `SELECT id, name, price, unit FROM products WHERE id=$1;`
	err = db.QueryRow(query, ProductId).Scan(&existingProduct.ID, &existingProduct.Name, &existingProduct.Price, &existingProduct.Unit)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Product"})
		}
		return
	}

	// Pengkondisian: jika setiap field tidak kosong, kita akan update datanya. Jika kosong, kita akan menggunakan data sebelumnya atau partial update
	if strings.TrimSpace(updatedProduct.Name) != "" {
		existingProduct.Name = updatedProduct.Name
	}
	if updatedProduct.Price != 0 { // Assuming zero is not a valid price, otherwise you need a different condition
		existingProduct.Price = updatedProduct.Price
	}
	if strings.TrimSpace(updatedProduct.Unit) != "" {
		existingProduct.Unit = updatedProduct.Unit
	}

	// Update data pelanggan di database
	updateQuery := `UPDATE products SET name=$1, price=$2, unit=$3 WHERE id=$4`
	_, err = db.Exec(updateQuery, existingProduct.Name, existingProduct.Price, existingProduct.Unit, ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Product"})
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "data": existingProduct})

	// Membuat respons dengan struktur yang diinginkan
	response := ResponseProduct{
		Message: "Product updated successfully",
		Data:    existingProduct,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteProductById(c *gin.Context) {
	id := c.Param("id")

	ProductId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product Id"})
		return
	}

	query := `DELETE FROM products WHERE id=$1;`
	_, err = db.Exec(query, ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete Product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"data":    "OK",
	})

}
