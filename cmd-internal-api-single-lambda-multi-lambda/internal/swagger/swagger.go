package swagger

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/jha-captech/user-microservice/internal/swagger/docs"
	"github.com/swaggo/http-swagger/v2"
)

type sLogger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

func RunSwagger(r *chi.Mux, logger sLogger, host string) {
	// docs
	docs.SwaggerInfo.Title = "User Microservice API"
	docs.SwaggerInfo.Description = "Sample Go API"
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/api"

	docs.SwaggerInfo.Schemes = []string{"http"}

	// handler
	baseURL := "http://" + host

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(baseURL+"/swagger/doc.json"),
	))

	logger.Info(fmt.Sprintf("Swagger URL: %s/swagger/index.html", baseURL))
}
