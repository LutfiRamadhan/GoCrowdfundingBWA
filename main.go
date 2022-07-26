package main

import (
	"BWA/auth"
	"BWA/campaign"
	"BWA/handler"
	"BWA/helper"
	"BWA/user"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// router := gin.Default()

	dsn := "host=127.0.0.1 user=postgres password=p4ssw0rd dbname=cfd_bwa port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)
	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images/")
	api := router.Group("api/v1")
	api.POST("/register", userHandler.RegisterUser)
	api.POST("/login", userHandler.LoginUser)
	api.POST("/validate_email", userHandler.CheckEmail)
	api.POST("/upload_avatar", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/campaign", authMiddleware(authService, userService), campaignHandler.ListCampaign)
	api.GET("/campaign/:slug", authMiddleware(authService, userService), campaignHandler.GetCampaign)
	api.POST("/campaign", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.POST("/campaign_images", authMiddleware(authService, userService), campaignHandler.CreateImage)
	api.GET("/export_campaign", campaignHandler.ExportCampaign)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer ") {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) != 2 {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := arrayToken[1]
		if tokenString == "" {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.FindByID(userID)
		if err != nil {
			response := helper.ResponseAPI("Unauthorized", http.StatusUnauthorized, "Error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set("currentUser", user)
		return
	}

}
