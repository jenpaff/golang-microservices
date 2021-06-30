package errors

import (
	"errors"
	"net/http"
)

type httpError struct {
	error  string
	status int
}

type ErrorResponse struct {
	ErrorID      string `json:"error_id"`
	ErrorMessage string `json:"error_message"`
}

var httpErrors = make(map[error]*httpError)

var InternalServerError = newHttpError("INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
var BadRequest = newHttpError("BAD_REQUEST", http.StatusBadRequest)
var UserNotFound = newHttpError("USER_NOT_FOUND", http.StatusNotFound)
var UserClientError = newHttpError("USER_CLIENT_ERROR", http.StatusInternalServerError)

func newHttpError(errorID string, status int) *httpError {
	error := &httpError{
		error:  errorID,
		status: status,
	}
	httpErrors[error] = error
	return error
}

func (e *httpError) Error() string {
	return e.error
}

func (e *httpError) toErrorResponse(details string) (int, ErrorResponse) {
	errorResponse := ErrorResponse{
		ErrorID:      e.error,
		ErrorMessage: details,
	}

	return e.status, errorResponse
}

func httpErrorFromError(wrappedError error) *httpError {
	baseError := errors.Unwrap(wrappedError)
	if httpError, ok := httpErrors[baseError]; ok {
		return httpError
	}
	return InternalServerError
}
