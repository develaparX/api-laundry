package handlers

import (
	"database/sql"
	"enigma-laundry/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseEmployee struct {
	Message string          `json:"message"`
	Data    models.Employee `json:"data"`
}

type ResponseListEmployee struct {
	Message string            `json:"message"`
	Data    []models.Employee `json:"data"`
}

func CreateEmployee(c *gin.Context) {
	var newEmployee models.Employee
	err := c.ShouldBind(&newEmployee)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//simpan data ke database
	query := "INSERT INTO Employees(name, phone_number, address) VALUES ($1,$2,$3) RETURNING id"

	err = db.QueryRow(query, newEmployee.Name, newEmployee.PhoneNumber, newEmployee.Address).Scan(&newEmployee.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new Employee"})
		return
	}

	response := ResponseEmployee{
		Message: "Employee created successfully",
		Data:    newEmployee,
	}

	c.JSON(http.StatusCreated, response)
}

func GetAllEmployees(c *gin.Context) {
	searchId := c.Param("id")

	query := "SELECT id, name, phone_number, address FROM Employees"

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

	var matchedEmployee []models.Employee

	for rows.Next() {
		var Employee models.Employee

		err := rows.Scan(&Employee.ID, &Employee.Name, &Employee.PhoneNumber, &Employee.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
			return
		}
		matchedEmployee = append(matchedEmployee, Employee)
	}

	response := ResponseListEmployee{
		Message: "Employee updated successfully",
		Data:    matchedEmployee,
	}

	if len(matchedEmployee) > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Employee not found",
		})
	}

}

func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	EmployeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee id"})
	}

	var updatedEmployee models.Employee
	if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	//mengambil data pelanggan yang ada di database berdasarkan ID
	var existingEmployee models.Employee
	query := `SELECT id, name, phone_number, address FROM Employees WHERE id=$1;`
	err = db.QueryRow(query, EmployeeId).Scan(&existingEmployee.ID, &existingEmployee.Name, &existingEmployee.PhoneNumber, &existingEmployee.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Employee"})
		}
		return
	}

	// Pengkondisian: jika setiap field tidak kosong, kita akan update datanya. Jika kosong, kita akan menggunakan data sebelumnya atau partial update
	if strings.TrimSpace(updatedEmployee.Name) != "" {
		existingEmployee.Name = updatedEmployee.Name
	}
	if strings.TrimSpace(updatedEmployee.PhoneNumber) != "" {
		existingEmployee.PhoneNumber = updatedEmployee.PhoneNumber
	}
	if strings.TrimSpace(updatedEmployee.Address) != "" {
		existingEmployee.Address = updatedEmployee.Address
	}

	// Update data pelanggan di database
	updateQuery := `UPDATE Employees SET name=$1, phone_number=$2, address=$3 WHERE id=$4`
	_, err = db.Exec(updateQuery, existingEmployee.Name, existingEmployee.PhoneNumber, existingEmployee.Address, EmployeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Employee"})
		return
	}
	// c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully", "data": existingEmployee})

	// Membuat respons dengan struktur yang diinginkan
	response := ResponseEmployee{
		Message: "Employee updated successfully",
		Data:    existingEmployee,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteEmployeeById(c *gin.Context) {
	id := c.Param("id")

	EmployeeId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Employee Id"})
		return
	}

	query := `DELETE FROM Employees WHERE id=$1;`
	_, err = db.Exec(query, EmployeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete Employee"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Employee deleted successfully",
		"data":    "OK",
	})

}
