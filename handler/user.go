package handler

import (
	"BWA/helper"
	"BWA/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	var input user.RegisterUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseAPI("Failed to register account!", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		fmt.Println("Error register user", err.Error())
		response := helper.ResponseAPI("Failed to register account!", http.StatusInternalServerError, "Error", nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	json := user.FormatUser(newUser, "tokentokentoken")

	response := helper.ResponseAPI("Account has been registered!", http.StatusOK, "Success", json)
	ctx.JSON(http.StatusOK, response)
}
