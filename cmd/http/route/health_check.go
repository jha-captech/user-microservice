package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheck(r *gin.RouterGroup, h Handler) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})
}
