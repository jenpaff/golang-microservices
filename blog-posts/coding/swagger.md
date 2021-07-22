# Swagger OpenAPI Documentation

## Motivation
Spring Boot almost spoils you with the ease of generating Swagger OpenAPI Documentation 
by adding a few annotations and configuration. 
With such an experience in mind, we've decided to switch our router implementation with [go-restful](https://github.com/emicklei/go-restful), 
a package for building REST-style Web Services using Google Go which comes with a Swagger Documentation built-in.

## Usage 
`go-restful` allows us to add documentation at our REST endpoint definition, e.g. 

```go
ws.Route(
    ws.GET("/users/{userName}").
        To(errors.ErrorHandler(controller.GetUser)).
        Doc("get users endpoint").
        Param(ws.PathParameter("userName", "name of the user").DataType("string")).
        Writes(common.User{}).
        Metadata(restfulspec.KeyOpenAPITags, tagsUser).
        Produces(restful.MIME_JSON).
        Returns(http.StatusOK, http.StatusText(http.StatusOK), common.User{}))
```

The Swagger Documentation for this endpoint looks something like this:

```json
"/users/{userName}": {
   "get": {
    "produces": [
     "application/json"
    ],
    "tags": [
     "user"
    ],
    "summary": "get users endpoint",
    "operationId": "func1",
    "parameters": [
     {
      "type": "string",
      "description": "name of the user",
      "name": "userName",
      "in": "path",
      "required": true
     }
    ],
    "responses": {
     "200": {
      "description": "OK",
      "schema": {
       "$ref": "#/definitions/common.User"
      }
     }
    }
   }
}
```

- [ ] add code to automatically display in the UI

## Code

As of now, `go-restful` automatically generates a `swagger.json` by configuring swagger doc for our router as follows: 

```go
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
```

The `swagger.json` can be retrieved by starting the container and calling 

```shell
curl localhost:12345/swagger.json
```

- [ ] add explanation of displaying in UI

## Further Resources
- [Display via swagger-ui](https://ribice.medium.com/serve-swaggerui-within-your-golang-application-5486748a5ed4) 
- [Serving in swagger-ui](https://gist.github.com/StevenACoffman/fe5f7774c750926210b642a0997479b0)