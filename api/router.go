package api

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/errors"
	"net/http"
)

var (
	tagsUser = []string{"user"}
)

func NewRouter(controller *Controller) http.Handler {
	wsContainer := restful.NewContainer()
	registerCorsFilter(wsContainer)

	ws := newService(controller)
	wsContainer.Add(ws)

	swaggerConfig := restfulspec.Config{
		WebServices: wsContainer.RegisteredWebServices(),
		APIPath:     "/swagger.json",
		PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
			swo.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Golang Service",
					Description: "A service that is written in Golang",
					Version:     "1.0.0",
				},
			}
		},
	}

	swaggerWs := restfulspec.NewOpenAPIService(swaggerConfig)
	wsContainer.Add(swaggerWs)

	return wsContainer
}

func newService(controller *Controller) *restful.WebService {
	ws := new(restful.WebService)
	ws.Route(
		ws.GET("/health").
			Doc("health endpoint").
			Writes(Health{}).
			To(controller.Health).
			Produces(restful.MIME_JSON).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), Health{}))

	ws.Route(
		ws.GET("/users/{userName}").
			To(errors.ErrorHandler(controller.GetUser)).
			Doc("get users endpoint").
			Param(ws.PathParameter("userName", "name of the user").DataType("string")).
			Writes(common.User{}).
			Metadata(restfulspec.KeyOpenAPITags, tagsUser).
			Produces(restful.MIME_JSON).
			Consumes(restful.MIME_JSON).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), common.User{}))

	ws.Route(
		ws.POST("/users").
			To(errors.ErrorHandler(controller.CreateUser)).
			Doc("create  users endpoint").
			Writes(common.User{}).
			Produces(restful.MIME_JSON).
			Reads(UserCreationRequest{}).
			Consumes(restful.MIME_JSON).
			Metadata(restfulspec.KeyOpenAPITags, tagsUser).
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
