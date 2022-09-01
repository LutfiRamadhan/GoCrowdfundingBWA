package handler

import (
	"BWA/campaign"
	"BWA/helper"
	"BWA/user"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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

	response := helper.ResponseAPI("Success create campaign", http.StatusOK, "Success", newCampaign)
	ctx.JSON(http.StatusBadRequest, response)
	return
}

func (h *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var input campaign.CampaignInput
	tmp_id := ctx.Param("id")
	id, err := strconv.Atoi(tmp_id)
	if id != 0 || err != nil {
		errorMessage := gin.H{"errors": "Wrong Format Parameter"}
		response := helper.ResponseAPI("Invalid request", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)
	input.UserId = currentUser.ID
	input.ID = id

	if err := ctx.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Failed to update campaign", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newCampaign, err := h.campaignService.UpdateCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseAPI("Success update campaign", http.StatusOK, "Success", newCampaign)
	ctx.JSON(http.StatusBadRequest, response)
	return
}

func (h *campaignHandler) ExportCampaign(ctx *gin.Context) {
	var input campaign.CampaignExportFilter
	if err := ctx.ShouldBindQuery(&input); err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	campaignData, err := h.campaignService.ExportCampaign(input)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	xls := excelize.NewFile()
	sheet1Name := "Campaign"
	xls.SetSheetName(xls.GetSheetName(0), sheet1Name)
	xls.SetCellValue(sheet1Name, "A1", "No")
	xls.SetCellValue(sheet1Name, "B1", "ID")
	xls.SetCellValue(sheet1Name, "C1", "Name")
	xls.SetCellValue(sheet1Name, "D1", "Short Description")
	xls.SetCellValue(sheet1Name, "E1", "Description")
	xls.SetCellValue(sheet1Name, "F1", "Goal Amount")
	xls.SetCellValue(sheet1Name, "G1", "Current Amount")
	xls.SetCellValue(sheet1Name, "H1", "Perks")
	xls.SetCellValue(sheet1Name, "I1", "Backer Count")
	xls.SetCellValue(sheet1Name, "J1", "slug")
	xls.SetCellValue(sheet1Name, "K1", "Created At")

	if err := xls.AutoFilter(sheet1Name, "A1", "L1", ""); err != nil {
		log.Fatal("Error: ", err.Error())
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	for i, data := range campaignData {
		createdAt, err := time.Parse(time.RFC822, data.CreatedAt.String())
		if err != nil {
			log.Fatal("Error: ", err.Error())
			createdAt = time.Now()
		}
		formattedCreatedAt := createdAt.Format("02/01/2006 15:04:05")

		xls.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), i+1)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), data.ID)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), data.Name)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("D%d", i+2), data.ShortDescription)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("E%d", i+2), data.Description)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("F%d", i+2), data.GoalAmount)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("G%d", i+2), data.CurrentAmount)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("H%d", i+2), data.Perks)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("I%d", i+2), data.BackerCount)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("J%d", i+2), data.Slug)
		xls.SetCellValue(sheet1Name, fmt.Sprintf("K%d", i+2), formattedCreatedAt)
	}
	var b bytes.Buffer
	if err := xls.Write(&b); err != nil {
		log.Fatal("Error: ", err.Error())
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Invalid request", http.StatusBadRequest, "Error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	downloadName := time.Now().UTC().Format("CAMPAIGN_02012006150405.xlsx")
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename="+downloadName)
	ctx.Data(http.StatusOK, "application/octet-stream", b.Bytes())
	return
}

func (h *campaignHandler) CreateImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		fmt.Println("Error upload campaign image", err.Error())
		response := helper.ResponseAPI("Error upload campaign image", http.StatusUnprocessableEntity, "Error", "Invalid input file")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var input campaign.CampaignImageInput
	if err := ctx.ShouldBind(&input); err != nil {
		errorMessage := gin.H{"errors": err}
		response := helper.ResponseAPI("Failed to create campaign image", http.StatusUnprocessableEntity, "Error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(user.User)
	path := fmt.Sprintf("./images/%d_%s", currentUser.ID, file.Filename)
	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		fmt.Println("Error upload campaign image", err.Error())
		response := helper.ResponseAPI("Error upload campaign image", http.StatusBadRequest, "Error", "Invalid input file")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	input.FileName = path
	input.UserId = currentUser.ID
	result, err := h.campaignService.CreateImage(input)
	if err != nil {
		fmt.Println("Error upload campaign image", err.Error())
		response := helper.ResponseAPI("Error upload campaign image", http.StatusBadRequest, "Error", "Bad Request!")
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseAPI("Successfully upload campaign image", http.StatusOK, "Success", result)
	ctx.JSON(http.StatusOK, response)
}
