package controller

import (
	"clean-laundry/service"
	"clean-laundry/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service service.ProductService
	rg      *gin.RouterGroup
}

func (p *ProductController) getAllHandler(ctx *gin.Context) {
	//default param itu string, jadi di ubah menjadi int langsung
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	size, err2 := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err2 != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	data, paging, err := p.service.FindAll(page, size)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	var listData []interface{}

	for _, pl := range data {
		listData = append(listData, pl)
	}

	utils.SendPagingResponse(ctx, "Success Get Data", listData, paging, http.StatusOK)
}

// route
func (p *ProductController) Route() {
	router := p.rg.Group("/products")
	{
		router.GET("/", p.getAllHandler)
	}
}

func NewProductController(service service.ProductService, rg *gin.RouterGroup) *ProductController {
	return &ProductController{
		service: service,
		rg:      rg,
	}
}
