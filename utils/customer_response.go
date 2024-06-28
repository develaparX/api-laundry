package utils

import (
	"clean-laundry/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func singleresponse
func SendSingleResponse(ctx *gin.Context, message string, data any, code int) {
	ctx.JSON(http.StatusOK, dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
		Data: data,
	})
}

func SendPagingResponse(ctx *gin.Context, message string, data []any, paging dto.Paging, code int) {
	ctx.JSON(http.StatusOK, dto.PagingResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendErrorResponse(ctx *gin.Context, message string, code int) {
	ctx.JSON(code, dto.SingleResponse{
		Status: dto.Status{
			Code:    code,
			Message: message,
		},
	})
}
