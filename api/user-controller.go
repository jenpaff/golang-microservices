package api

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/errors"
	"github.com/jenpaff/golang-microservices/featuretoggles"
	"github.com/jenpaff/golang-microservices/validation"
	"io/ioutil"
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

func (c *Controller) CreateUser(req *restful.Request, resp *restful.Response) error {

	log.Info("save user endpoint was invoked")

	var err error

	bytes, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return fmt.Errorf("could read request body: %w", err)
	}

	var creationRequest UserCreationRequest
	err = json.Unmarshal(bytes, &creationRequest)
	if err != nil {
		return fmt.Errorf("could not unmarshal the user request: %w", err)
	}

	err = c.validator.Struct(creationRequest)
	if err != nil {
		return validation.GetValidationError(err)
	}

	var createdUser *common.User

	ft := featuretoggles.NewFeatureToggles(&c.Cfg, req)
	if ft.IsEnabled("enableNewFeature") {
		createdUser, err = c.userService.CreateUserWithNewFeature(req.Request.Context(), creationRequest.UserName, creationRequest.Email, creationRequest.PhoneNumber)
		if err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}
	} else {
		createdUser, err = c.userService.CreateUser(req.Request.Context(), creationRequest.UserName, creationRequest.Email, creationRequest.PhoneNumber)
		if err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}
	}

	err = resp.WriteEntity(createdUser)
	if err != nil {
		log.Errorf("could not write response: %s", err.Error())
	}

	return nil

}
