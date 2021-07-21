package errors

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/go-playground/log"
)

type ErrorRouteFunction func(req *restful.Request, resp *restful.Response) error

func logAndReturnHttpError(err error) (int, ErrorResponse) {
	log.Error(err.Error())

	httpError := httpErrorFromError(err)
	return httpError.toErrorResponse(err.Error())
}

func ErrorHandler(errorRouteFunction ErrorRouteFunction) restful.RouteFunction {
	return func(req *restful.Request, res *restful.Response) {
		err := errorRouteFunction(req, res)
		if err != nil {
			err2 := res.WriteHeaderAndEntity(logAndReturnHttpError(err))
			if err2 != nil {
				log.Errorf("could not write error response: %s", err.Error())
			}
		}
	}
}
