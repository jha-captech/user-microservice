package main

import (
	"fmt"
	"log"

	"user-microservice/cmd/http/route"
	"user-microservice/internal/database"
	"user-microservice/internal/user"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	RunHTTP()
}

func RunHTTP() {
	config := mustNewConfiguration()

	db := database.MustNewDatabase(
		postgres.Open(
			fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
				config.Database.Host,
				config.Database.User,
				config.Database.Password,
				config.Database.Name,
				config.Database.Port,
			),
		),
		gorm.Config{},
		config.Database.ConnectionRetry,
	)
	us := user.NewService(db)

	h := route.NewHandler(us)

	r := gin.Default()

	route.SetUpRoutes(r, h)

	log.Fatal(r.Run(config.HTTP.Port))
}
