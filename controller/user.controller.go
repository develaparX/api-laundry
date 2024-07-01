package controller

import (
	"clean-laundry/model"
	"clean-laundry/model/dto"
	"clean-laundry/service"
	"clean-laundry/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
	rg      *gin.RouterGroup
}

func (u *UserController) loginHandler(ctx *gin.Context) {
	var payload dto.LoginDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.SendErrorResponse(ctx, "Failed to parsing payload", http.StatusBadRequest)
		return
	}

	response, errors := u.service.Login(payload)
	if errors != nil {
		utils.SendErrorResponse(ctx, errors.Error(), http.StatusInternalServerError)
		return
	}
	utils.SendSingleResponse(ctx, "Success Login", response, http.StatusOK)
}

func (u *UserController) registerHandler(ctx *gin.Context) {
	var payload model.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.SendErrorResponse(ctx, "Failed to Parsing Payload", http.StatusBadRequest)
	}

	data, errors := u.service.CreateNew(payload)
	if errors != nil {
		utils.SendErrorResponse(ctx, errors.Error(), http.StatusInternalServerError)
	}
	utils.SendSingleResponse(ctx, "Success Create New User", data, http.StatusOK)
}

func (u *UserController) findByUsernameHandler(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := u.service.FindByUsername(username)
	if err != nil {
		utils.SendErrorResponse(ctx, err.Error(), http.StatusBadRequest)
	}
	utils.SendSingleResponse(ctx, "user retrieved successfully", user, http.StatusOK)
}

func (u *UserController) Route() {
	router := u.rg.Group("/users")
	router.POST("/register", u.registerHandler)
	router.GET("/:username", u.findByUsernameHandler)
	router.GET("/login", u.loginHandler)
}

func NewUserController(uS service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{
		service: uS,
		rg:      rg,
	}
}
