package handler

import (
	"BWA/campaign"
	"BWA/helper"
	"BWA/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
}

func NewCampaignHandler(campaignService campaign.Service) *campaignHandler {
	return &campaignHandler{campaignService}
}

func (h *campaignHandler) ListCampaign(ctx *gin.Context) {
	var input campaign.CampaignInput

	if filterName := ctx.Query("name"); filterName != "" {
		input.Name = filterName
	}

	if filterUserId := ctx.Query("user_id"); filterUserId != "" {
		userId, err := strconv.Atoi(filterUserId)
		if err != nil {
			errorMessage := gin.H{"errors": err}
			response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		input.UserId = userId
	}

	campaignData, err := h.campaignService.FindCampaigns(input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseAPI("Login Success!", http.StatusOK, "Success", campaignData)
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(ctx *gin.Context) {
	var input campaign.CampaignInput
	if slug := ctx.Param("slug"); slug == "" {
		errorMessage := gin.H{"errors": "Wrong Format Parameter"}
		response := helper.ResponseAPI("Invalid request", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.Slug = ctx.Param("slug")
	campaignData, err := h.campaignService.FindCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseAPI("Detail Campaign", http.StatusOK, "Success", campaignData)
	ctx.JSON(http.StatusBadRequest, response)
	return
}

func (h *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var input campaign.CampaignInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Failed to create campaign", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)
	userId := currentUser.ID

	newCampaign, err := h.campaignService.CreateCampaign(input, userId)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseAPI("Succes create campaign", http.StatusOK, "Success", newCampaign)
	ctx.JSON(http.StatusBadRequest, response)
	return
}
