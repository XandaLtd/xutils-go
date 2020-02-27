package xerrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RestErr interface {
	Error() bool
	StatusCode() int
	Message() string
}

type restErr struct {
	ErrError      bool   `json:"error"`
	ErrStatusCode int    `json:"status_code"`
	ErrMessage    string `json:"message"`
}

func (e restErr) Error() bool {
	return e.ErrError
}

func (e restErr) StatusCode() int {
	return e.ErrStatusCode
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func NewRestError(status int, message string) RestErr {
	return restErr{
		ErrError:      true,
		ErrStatusCode: status,
		ErrMessage:    message,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid error json response")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrError:      true,
		ErrStatusCode: http.StatusBadRequest,
		ErrMessage:    message,
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrError:      true,
		ErrStatusCode: http.StatusNotFound,
		ErrMessage:    message,
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrError:      true,
		ErrStatusCode: http.StatusUnauthorized,
		ErrMessage:    message,
	}
}

func NewInternalServerError(message string) RestErr {
	return restErr{
		ErrError:      true,
		ErrStatusCode: http.StatusInternalServerError,
		ErrMessage:    message,
	}
}
