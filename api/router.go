package api

import (
	"github.com/emicklei/go-restful"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/errors"
	"net/http"
)

func NewRouter(controller *Controller) http.Handler {
	wsContainer := restful.NewContainer()
	registerCorsFilter(wsContainer)
	ws := newService(controller)
	wsContainer.Add(ws)
	return wsContainer
}

func newService(controller *Controller) *restful.WebService {
	ws := new(restful.WebService)
	ws.Route(
		ws.GET("/health").
			To(controller.Health).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), Health{}))
	ws.Route(
		ws.GET("/users/{userName}").
			To(errors.ErrorHandler(controller.GetUser)).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), common.User{}))
	ws.Route(
		ws.POST("/users").
			To(errors.ErrorHandler(controller.CreateUser)).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), common.User{}))
	return ws
}

func registerCorsFilter(wsContainer *restful.Container) {
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
}
