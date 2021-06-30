package api

import (
	"github.com/jenpaff/golang-microservices/config"
	"github.com/jenpaff/golang-microservices/users"
)

type Controller struct {
	cfg         config.Config
	userService users.Service
}

func NewController(cfg config.Config) *Controller {
	return &Controller{cfg: cfg}
}
