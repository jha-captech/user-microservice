package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jha-captech/user-microservice/internal/models"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

type ResponseID struct {
	ObjectID int `json:"object_id"`
}

type ResponseOneUser struct {
	User models.User `json:"user"`
}

type ResponseAllUsers struct {
	Users []models.User `json:"users"`
}
type ResponseError struct {
	Error string `json:"error"`
}

func (h *Handler) returnErr(msg string, err error, status int) (events.APIGatewayProxyResponse, error) {
	h.logger.Error(msg, "err", err)

	respBody, _ := json.Marshal(ResponseError{Error: msg})

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(respBody),
	}, nil
}

func (h *Handler) returnJSON(statusCode int, data any) (events.APIGatewayProxyResponse, error) {
	JSONData, err := json.Marshal(data)
	if err != nil {
		return h.returnErr("Error marshaling return", err, http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(JSONData),
	}, nil
}
