package route

import (
	"log/slog"
	"net/http"
	"strconv"

	"user-microservice/internal/database/entity"

	"github.com/gin-gonic/gin"
)

type responseOneUser struct {
	User entity.User `json:"user"`
}

type responseAllUsers struct {
	Users []entity.User `json:"users"`
}

func userRoutes(r *gin.RouterGroup, h Handler) {
	// list all users
	r.GET("/", func(c *gin.Context) {
		// get values from db
		users, err := h.userService.List()
		if err != nil {
			slog.Error("error getting all locations", "error", err)
			c.JSON(http.StatusInternalServerError, responseError{Error: "Error retrieving data"})
			return
		}

		// return response
		c.JSON(http.StatusOK, responseAllUsers{Users: users})
	})

	// fetch a user by ID
	r.GET("/:ID", func(c *gin.Context) {
		// get and validate ID
		idString := c.Param("ID")
		ID, err := strconv.Atoi(idString)
		if err != nil {
			slog.Error("error getting ID", "error", err)
			c.JSON(http.StatusBadRequest, responseError{Error: "Not a valid ID"})
			return
		}

		// get values from db
		user, err := h.userService.Fetch(ID)
		if err != nil {
			slog.Error("error getting all locations", "error", err)
			c.JSON(http.StatusInternalServerError, responseError{Error: "Error retrieving data"})
			return
		}

		// return response
		c.JSON(http.StatusOK, responseOneUser{User: user})
	})
}
