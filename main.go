package main

import (
	"BWA/handler"
	"BWA/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// router := gin.Default()

	dsn := "host=172.25.2 user=postgres password=p4ssw0rd dbname=cfd_bwa port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	UserRepository := user.NewRepository(db)
	userService := user.NewService(UserRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("api/v1")
	api.POST("/user", userHandler.RegisterUser)

	router.Run()
}
