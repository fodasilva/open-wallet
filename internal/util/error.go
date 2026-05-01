package util

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
