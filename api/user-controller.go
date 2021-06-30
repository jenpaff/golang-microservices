package api

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/errors"
	"strings"
)

func (c *Controller) GetUser(req *restful.Request, resp *restful.Response) error {

	log.Info("user endpoint was invoked")

	userName := req.PathParameter("userName")
	if strings.TrimSpace(userName) == "" {
		return fmt.Errorf("you must provide a valid username: %w", errors.BadRequest)
	}

	user, err := c.userService.GetUser(req.Request.Context(), userName)
	if err != nil {
		return fmt.Errorf("error retrieving user with username %s: %w", userName, err)
	}
	err = resp.WriteEntity(user)
	if err != nil {
		return fmt.Errorf("error writing response: %w", err)
	}
	log.Info("health endpoint ran successfully")

	return nil
}
