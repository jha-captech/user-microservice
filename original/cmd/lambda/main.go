package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"gorm.io/driver/postgres"

	"user-microservice/cmd/lambda/handler"
	"user-microservice/internal/database"
	"user-microservice/internal/user"
)

func main() {
	lambda.Start(run())
}

func run() handler.APIGatewayHandler {
	config := mustNewConfiguration()

	logger := newLogger()

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
		database.WithLogger(logger),
		database.WithRetryCount(5),
		database.WithAutoMigrate(true),
	)

	us := user.NewService(db)

	h := handler.New(us, logger)

	return handler.Run(h)
}
