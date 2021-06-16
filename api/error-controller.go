package api

import (
	"fmt"
	"github.com/emicklei/go-restful"
)

func (c *Controller) Error(req *restful.Request, resp *restful.Response) error {
	return fmt.Errorf("could not process incoming request %w", ErrBadRequest)
}
