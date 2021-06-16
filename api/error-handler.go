package api

import (
	"errors"
	. "github.com/emicklei/go-restful"
	"github.com/go-playground/log"
	"net/http"
)

type HttpError interface {
	Error() string
	ToErrorResponse(string) (int, ErrorResponse)
}

type httpError struct {
	error  string
	status int
}

type ErrorResponse struct {
	ErrorID      string `json:"error_id"`
	ErrorMessage string `json:"error_message"`
}

var httpErrors = make(map[HttpError]bool)

var ErrInternalServerError = newHttpError("INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
var ErrBadRequest = newHttpError("BAD_REQUEST", http.StatusBadRequest)
var ErrUserNotFound = newHttpError("USER_NOT_FOUND", http.StatusNotFound)

func newHttpError(errorID string, status int) HttpError {
	error := &httpError{
		error:  errorID,
		status: status,
	}
	httpErrors[error] = true
	return error
}

func (e *httpError) Error() string {
	return e.error
}

func (e *httpError) ToErrorResponse(details string) (int, ErrorResponse) {
	errorResponse := ErrorResponse{
		ErrorID:      e.error,
		ErrorMessage: details,
	}

	return e.status, errorResponse
}

func LogAndReturnHttpError(err error) (int, ErrorResponse) {
	log.Error(err.Error())

	httpError, _ := httpErrorFromError(err)
	return httpError.ToErrorResponse(err.Error())
}

func httpErrorFromError(wrappedError error) (HttpError, bool) {
	for httpError := range httpErrors {
		if errors.Is(wrappedError, httpError) {
			return httpError, true
		}
	}
	return ErrInternalServerError, false
}

// example usage
func ErrorHandler(req *Request, res *Response) {
	var err error // = some function call
	if err != nil {
		err2 := res.WriteHeaderAndEntity(LogAndReturnHttpError(err))
		if err2 != nil {
			log.Errorf("could not write error response: %s", err.Error())
		}
	}
}
