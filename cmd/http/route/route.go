package route

import (
	"github.com/gin-gonic/gin"

	"user-microservice/internal/database/entity"
)

type responseMessage struct {
	Message string `json:"message"`
}

type responseID struct {
	ObjectID int `json:"object_id"`
}

type responseError struct {
	Error string `json:"error"`
}

type userService interface {
	List() ([]entity.User, error)
	Fetch(ID int) (entity.User, error)
}

type Handler struct {
	userService userService
}

func NewHandler(userService userService) Handler {
	return Handler{userService: userService}
}

func SetUpRoutes(r *gin.Engine, h Handler) {
	healthCheck(r.Group("/health-check"), h)

	userRoutes(r.Group("/user"), h)

	notFound(r)
}
