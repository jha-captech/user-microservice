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

func (h *Handler) returnJSON(statusCode int, data any) (events.APIGatewayProxyResponse, error) {
	JSONData, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Error marshaling return", "err", err, "data", data)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: statusCode, Body: string(JSONData)}, nil
}
