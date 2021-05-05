package api

import (
	"github.com/emicklei/go-restful"
	log "github.com/mgutz/logxi/v1"
)

type Health struct {
	status string
}

func (c *Controller) Health(_ *restful.Request, resp *restful.Response) {
	health := Health{status: "up"}
	err := resp.WriteEntity(health)
	if err != nil {
		log.Warn("service is down, cannot write health status")
	}
}
