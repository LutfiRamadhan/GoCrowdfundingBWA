package handler

import (
	"BWA/auth"
	"BWA/helper"
	"BWA/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(ctx *gin.Context) {
	var input user.RegisterUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseAPI("Failed to register account!", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		fmt.Println("Error register user", err.Error())
		response := helper.ResponseAPI("Failed to register account!", http.StatusBadRequest, "Error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		fmt.Println("Error getting token", err.Error())
		response := helper.ResponseAPI("Failed getting token, please re-login!", http.StatusBadRequest, "Error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	json := user.FormatUser(newUser, token)

	response := helper.ResponseAPI("Account has been registered!", http.StatusOK, "Success", json)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) LoginUser(ctx *gin.Context) {
	var input user.LoginUserInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseAPI("Not Valid!", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	userData, err := h.userService.LoginUser(input)
	if err != nil {
		fmt.Println("Error login user", err.Error())
		response := helper.ResponseAPI(err.Error(), http.StatusBadRequest, "Error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(userData.ID)
	if err != nil {
		fmt.Println("Error getting token", err.Error())
		response := helper.ResponseAPI("Failed getting token, please re-login!", http.StatusBadRequest, "Error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	json := user.FormatUser(userData, token)

	response := helper.ResponseAPI("Login Success!", http.StatusOK, "Success", json)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmail(ctx *gin.Context) {
	var input user.ValidateEmailInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseAPI("Email not Valid", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	status, err := h.userService.ValidateEmail(input)
	if err != nil {
		fmt.Println("Error validate email", err.Error())
		response := helper.ResponseAPI(err.Error(), http.StatusBadRequest, "Error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	json := gin.H{
		"email": status,
	}
	response := helper.ResponseAPI("Success check email!", http.StatusOK, "Success", json)
	ctx.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(ctx *gin.Context) {
	data_res := gin.H{"is_uploaded": false}
	file, err := ctx.FormFile("avatar")
	if err != nil {
		fmt.Println("Error upload avatar", err.Error())
		response := helper.ResponseAPI("Error upload avatar", http.StatusUnprocessableEntity, "Error", data_res)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	currentUser := ctx.MustGet("currentUser").(user.User)
	path := fmt.Sprintf("./images/%d_%s", currentUser.ID, file.Filename)
	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		fmt.Println("Error upload avatar", err.Error())
		response := helper.ResponseAPI("Error upload avatar", http.StatusBadRequest, "Error", data_res)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.userService.SaveAvatar(currentUser.ID, path)
	if err != nil {
		fmt.Println("Error upload avatar", err.Error())
		response := helper.ResponseAPI("Error upload avatar", http.StatusBadRequest, "Error", data_res)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	data_res = gin.H{"is_uploaded": true}
	response := helper.ResponseAPI("Successfully upload avatar", http.StatusOK, "Success", data_res)
	ctx.JSON(http.StatusOK, response)
}
