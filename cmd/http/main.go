package main

import (
	"fmt"
	"log"
	"net/http"

	"user-microservice/cmd/http/route"
	"user-microservice/internal/database"
	"user-microservice/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"gorm.io/driver/postgres"
)

func main() {
	RunHTTP()
}

func RunHTTP() {
	config := mustNewConfiguration()

	logger := newLogger(false)

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

	h := route.NewHandler(us, logger)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	route.SetUpRoutes(r, h)

	log.Fatal(http.ListenAndServe(config.HTTP.Port, r))
}
