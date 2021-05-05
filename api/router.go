package api

import (
	"github.com/emicklei/go-restful"
	"github.com/rs/cors"
	"net/http"
)

func NewRouter(controller *Controller) http.Handler {
	ws := new(restful.WebService)
	ws.GET("/health").To(controller.Health)

	corsRouter := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
	})

	restfulContainer := restful.NewContainer()
	restfulContainer.Add(ws)

	return corsRouter.Handler(restfulContainer.ServeMux)
}
