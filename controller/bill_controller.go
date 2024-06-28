package controller

import (
	"clean-laundry/model/dto"
	"clean-laundry/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// controller tidak perlu bikin interface, tapi langsung struct aja
type BillController struct {
	service service.BillService
	rg      *gin.RouterGroup
}

// membuat handler yang akan di panggil di route
func (b *BillController) createHandler(ctx *gin.Context) {
	var payload dto.BillRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response, errs := b.service.CreateNewBill(payload)
	if errs != nil {
		ctx.JSON(http.StatusBadRequest, errs.Error())
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

// function route
func (b *BillController) Route() {
	group := b.rg.Group("/transactions")
	{
		group.POST("/", b.createHandler)
	}
}

// constructor untuk diakses dari luar
func NewBillController(service service.BillService, rg *gin.RouterGroup) *BillController {
	return &BillController{service: service, rg: rg}
}
