package errors

import (
	. "github.com/emicklei/go-restful"
	"github.com/go-playground/log"
)

type ErrorRouteFunction func(req *Request, resp *Response) error

func logAndReturnHttpError(err error) (int, ErrorResponse) {
	log.Error(err.Error())

	httpError := httpErrorFromError(err)
	return httpError.toErrorResponse(err.Error())
}

func ErrorHandler(errorRouteFunction ErrorRouteFunction) RouteFunction {
	return func(req *Request, res *Response) {
		err := errorRouteFunction(req, res)
		if err != nil {
			err2 := res.WriteHeaderAndEntity(logAndReturnHttpError(err))
			if err2 != nil {
				log.Errorf("could not write error response: %s", err.Error())
			}
		}
	}
}
