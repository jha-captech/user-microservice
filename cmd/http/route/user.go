package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"user-microservice/internal/database/entity"
)

type responseOneUser struct {
	User entity.User `json:"user"`
}

type responseAllUsers struct {
	Users []entity.User `json:"users"`
}

func userRoutes(h Handler) func(r chi.Router) {
	return func(r chi.Router) {
		// @Summary		List all users
		// @Description	List all users
		// @Tags		users
		// @Accept		json
		// @Produce		json
		// @Success		200		{object}	route.responseAllUsers
		// @Failure		500		{object}	route.responseError
		// @Router		/user	[GET]
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// get values from db
			users, err := h.userService.List()
			if err != nil {
				h.logger.Error("error getting all locations", "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error retrieving data"},
				)
				return
			}

			// return response
			encodeResponse(w, http.StatusOK, responseAllUsers{Users: users})
		})

		// @Summary		Fetch a user by ID
		// @Description	Fetch a user by ID
		// @Tags		user
		// @Accept		json
		// @Produce		json
		// @Param		id			path		int	true				"User ID"
		// @Success		200			{object}	route.responseOneUser
		// @Failure		400			{object}	route.responseError
		// @Failure		500			{object}	route.responseError
		// @Router		/user/{ID}	[GET]
		r.Get("/{ID}", func(w http.ResponseWriter, r *http.Request) {
			// get and validate ID
			idString := chi.URLParam(r, "ID")
			ID, err := strconv.Atoi(idString)
			if err != nil {
				h.logger.Error("error getting ID", "error", err)
				encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
				return
			}

			// get values from db
			user, err := h.userService.Fetch(ID)
			if err != nil {
				h.logger.Error("error getting locations buy ID", "ID", ID, "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error retrieving data"},
				)
				return
			}

			// return response
			encodeResponse(w, http.StatusOK, responseOneUser{User: user})
		})

		// @Summary		Update a user by ID
		// @Description	Update a user by ID
		// @Tags		user
		// @Accept		json
		// @Produce		json
		// @Param		id			path		int	true				"User ID"
		// @Param		user		body		entity.User	true		"User Object"
		// @Success		200			{object}	route.responseOneUser
		// @Failure		500			{object}	route.responseError
		// @Failure		422			{object}	route.responseError
		// @Router		/user/{ID}	[PUT]
		r.Put("/{ID}", func(w http.ResponseWriter, r *http.Request) {
			// get and validate ID
			idString := chi.URLParam(r, "ID")
			ID, err := strconv.Atoi(idString)
			if err != nil {
				h.logger.Error("error getting ID", "error", err)
				encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
				return
			}

			// get and validate body as object
			inputUser, err := decodeToStruct[entity.User](r)
			if err != nil {
				h.logger.Error("BodyParser error", "error", err)
				encodeResponse(
					w,
					http.StatusBadRequest,
					responseError{Error: "missing values or malformed body"},
				)
				return
			}

			// update object in database
			user, err := h.userService.Update(ID, inputUser)
			if err != nil {
				h.logger.Error("error updating object in db", "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error updating data"},
				)
				return
			}

			// return response
			encodeResponse(w, http.StatusOK, responseOneUser{User: user})
		})

		// @Summary		Create a user
		// @Description	Create a user
		// @Tags		user
		// @Accept		json
		// @Produce		json
		// @Param		user		body		entity.User	true	"User Object"
		// @Success		201			{object}	route.responseID
		// @Failure		422			{object}	route.responseError
		// @Failure		500			{object}	route.responseError
		// @Failure		409			{object}	route.responseError
		// @Router		/user		[POST]
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			// get and validate body as object
			inputUser, err := decodeToStruct[entity.User](r)
			if err != nil {
				h.logger.Error("BodyParser error", "error", err)
				encodeResponse(
					w,
					http.StatusBadRequest,
					responseError{Error: "missing values or malformed body"},
				)
				return
			}

			// create object in database
			id, err := h.userService.Create(inputUser)
			if err != nil {
				h.logger.Error("error creating object to db", "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error creating object"},
				)
				return
			}

			// return response
			encodeResponse(w, http.StatusOK, responseID{ObjectID: id})
		})

		// @Summary		Delete a user by ID
		// @Description	Delete a user by ID
		// @Tags		user
		// @Accept		json
		// @Produce		json
		// @Param		id			path		int	true				"User ID"
		// @Success		202			{object}	route.responseMessage
		// @Failure		500			{object}	route.responseError
		// @Failure		404			{object}	route.responseError
		// @Router		/user/{ID}	[DELETE]
		r.Delete("/{ID}", func(w http.ResponseWriter, r *http.Request) {
			// get and validate ID
			idString := chi.URLParam(r, "ID")
			ID, err := strconv.Atoi(idString)
			if err != nil {
				h.logger.Error("error getting ID", "error", err)
				encodeResponse(w, http.StatusBadRequest, responseError{Error: "Not a valid ID"})
				return
			}

			// check that object exists
			user, err := h.userService.Fetch(ID)
			if err != nil {
				h.logger.Error("error getting object by ID", "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error validating object"},
				)
				return
			}
			if user.ID == 0 {
				encodeResponse(w, http.StatusBadRequest, responseError{Error: "Object does not exist"})
				return
			}

			// delete user
			if err = h.userService.Delete(ID); err != nil {
				h.logger.Error("error deleting object by ID", "ID", ID, "error", err)
				encodeResponse(
					w,
					http.StatusInternalServerError,
					responseError{Error: "Error deleting object."},
				)
				return
			}

			// return message
			encodeResponse(w, http.StatusOK, responseMessage{Message: "object successful deleted"})
		})
	}
}
