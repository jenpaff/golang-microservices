package api

import "github.com/jenpaff/golang-microservices/users"

type Controller struct {
	userService users.Service
}

func NewController() *Controller {
	return &Controller{}
}
