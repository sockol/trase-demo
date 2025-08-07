package handlers

import "fmt"

type HTTPError struct {
	Code    int
	Message error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Error Code %d: %s", e.Code, e.Message)
}

func NewHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: err,
	}
}
