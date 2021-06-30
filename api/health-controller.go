package api

import (
	"github.com/emicklei/go-restful"
	"github.com/go-playground/log"
)

type Health struct {
	Status string `json:"status"`
	Name   string `json:"name"`
}

func (c *Controller) Health(_ *restful.Request, resp *restful.Response) {
	log.Info("health endpoint was invoked")
	health := &Health{Status: "up", Name: c.Cfg.Name}
	err := resp.WriteEntity(health)
	if err != nil {
		log.Warn("service is down, cannot write health status")
	}
	log.Info("health endpoint ran successfully")
}
