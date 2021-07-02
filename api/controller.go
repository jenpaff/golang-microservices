package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/users"
)

type Controller struct {
	Cfg         config.Config
	userService users.Service
	validator   *validator.Validate
}

func NewController(cfg config.Config, userService users.Service, validator *validator.Validate) *Controller {
	return &Controller{Cfg: cfg, userService: userService, validator: validator}
}
