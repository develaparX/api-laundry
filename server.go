package main

import (
	"clean-laundry/config"
	"clean-laundry/controller"
	"clean-laundry/repository"
	"clean-laundry/service"
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// server ini menghubungkan semua komponen
type Server struct {
	bS      service.BillService
	cS      service.CustomerService
	pS      service.ProductService
	uS      service.UserService
	engine  *gin.Engine //untuk start engine gin
	portApp string
}

// method untuk memanggil route yang di controller
func (s *Server) initiateRoute() {
	//bisa menambah grouping lagi disini
	routerGroup := s.engine.Group("/api/v1")
	controller.NewBillController(s.bS, routerGroup).Route()
	controller.NewProductController(s.pS, routerGroup).Route()
}

// func untuk running
func (s *Server) Start() {
	s.initiateRoute()
	s.engine.Run(s.portApp)
}

// constructur, agar dipanggil main.go
func NewServer() *Server {
	//memanggil hasil config .env
	co, _ := config.NewConfig()

	//melakukan koneksi database
	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)
	if err != nil {
		log.Fatal(err)
	}

	portApp := co.AppPort
	billRepo := repository.NewBillRepository(db)
	custRepo := repository.NewCustomerRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	custService := service.NewCustomerService(custRepo)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	billService := service.NewBillService(billRepo, userService, productService, custService)

	//menginject repo ke service
	return &Server{
		bS:      billService,
		cS:      custService,
		pS:      productService,
		uS:      userService,
		portApp: portApp,
		engine:  gin.Default(),
	}
}
