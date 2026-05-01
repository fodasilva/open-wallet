package httputil

import (
	"errors"
	"net/http"
)

var StatusMessages = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	429: "too many requests",
	500: "internal server error",
}

type HTTPErrorData struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type HTTPError struct {
	StatusCode int           `json:"status"`
	ErrorData  HTTPErrorData `json:"error"`
}

func (e *HTTPError) Error() string {
	return e.ErrorData.Message
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		ErrorData: HTTPErrorData{
			Type:    StatusMessages[statusCode],
			Message: message,
		},
	}
}

func GetApiErr(err error) *HTTPError {
	var apiErr *HTTPError

	if errors.As(err, &apiErr) {
		return apiErr
	} else {
		return NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}
}
