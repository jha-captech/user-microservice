package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseNotFound struct {
	Message string `json:"message"`
}

func notFound(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(
			http.StatusNotFound,
			responseNotFound{Message: "Page not found"},
		)
	})
}
