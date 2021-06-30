package api

import (
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/users"
)

type Controller struct {
	Cfg         config.Config
	userService users.Service
}

func NewController(cfg config.Config, userService users.Service) *Controller {
	return &Controller{Cfg: cfg, userService: userService}
}
