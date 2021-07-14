# Error Handling in REST Services

## Motivation

Error handling in Go can be quite painful. We want to share some best practices on how you keep your controllers clean and implement a centralized error handling avoiding redundancy.

Our requirement for error handling in REST services:

1. We want to be able to "throw" an error in every part of the code base
1. The error should be able to bubble up all the way to the controller
1. We want to be able to handle an error only once
1. Any error should be logged in a central error handler
1. Any error should possibly be wrapped in a custom error type that gets resolved to a corresponding HTTP status
1. Adding new error types bound to a specific HTTP status should be easy

## Custom Error Types in Go

All implementation related to error handling can be found in the [`api/error-handler.go` file](../api/error-handler.go)

### Implementing your own Custom Error Type

The built-in Go error interface is quite simple:

```Go
type error interface {
    Error() string
}
```

You can build your own custom error (e.g. `HttpError`) types by simply implementing this interface:

```Go
type HttpError interface {
	Error() string
	ToErrorResponse(string) (int, ErrorResponse)
}

func (e *httpError) Error() string {
    return e.error
}
```

Our implementation of this interfaces takes an error ID, e.g. `INVALID_REQUEST_BODY` and a HTTP status code, e.g. `http.StatusBadRequest`:

```Go
type httpError struct {
	error  string
	status int
}

func (e *httpError) Error() string {
    return e.error
}
```

A constructor method allows to create arbitrary custom errors of our `HttpError` interface:

```Go
var httpErrors = make(map[HttpError]bool)

func newHttpError(errorID string, status int) HttpError {
	error := &httpError{
		error:  errorID,
		status: status,
	}
	httpErrors[error] = true
	return error
}
```

You might be wondering about the `httpErrors[error] = true` in this constructor, we will come to this in a second.

We keep this method private on purpose, since we only want to export the custom errors we created with this method:

```Go
var ErrInternalServerError = newHttpError("INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
var ErrBadRequest = newHttpError("BAD_REQUEST", http.StatusBadRequest)
```

### Using Custom Error Types in your Code

A very simple usage of the custom error types is to wrap your error into such a custom error:

```Go
fmt.Errorf("could not process request %w", ErrBadRequest)
```

As you can see, you don't have to care about the associated HTTP status code or an error ID that will be shown to the client.

## Handling Errors on Router Level

The idea now is to let an error "bubble up" from the place in your code where it happens and pass it on to the Controller and finally to the Router. Therefore, we add an error handler that implements the `RouteFunction` interface for the restful router:

```Go
type ErrorRouteFunction func(req *Request, resp *Response) error

func ErrorHandler(errorRouteFunction ErrorRouteFunction) RouteFunction {
    return func(req *Request, res *Response) {
        err := errorRouteFunction(req, res)
        if err != nil {
            err2 := res.WriteHeaderAndEntity(LogAndReturnHttpError(err))
            if err2 != nil {
                log.Errorf("could not write error response: %s", err.Error())
            }
        }
    }
}
```

The error handler acts as a proxy around the controller method. It passes on the `req *http.Request` and `resp *http.Response` to the actual controller, catches the error and handles it, if it should occur.

You might be wondering about the `err2 := res.WriteHeaderAndEntity(LogAndReturnHttpError(err))`, be patient as this will be explained in a second!

We now expect our controller methods to implement the newly created `ErrorRouteFunction` interface, like this:

```Go
func (c *Controller)  ControllerMethod(req *restful.Request, resp *restful.Response) error {
	// ...
}
```

We can register the error handler in our router, since it implements the `RouteFunction` interface:

```Go
ws := new(restful.WebService)
ws.Route(
    ws.GET("/error").
        To(ErrorHandler(controller.Error)).
        Produces(restful.MIME_JSON).
        Returns(http.StatusOK, http.StatusText(http.StatusOK), nil))
```

### Dispatching Errors, Writing HTTP Response and Logging

The last missing piece is now the dispatching of the errors caught in the error handler, returning a proper error message to the (HTTP) client and logging the error.

We therefore add another function to our error handler that does all this:

```Go
func LogAndReturnHttpError(err error) (int, ErrorResponse) {
	log.Error(err.Error())

	httpError, _ := httpErrorFromError(err)
	return httpError.ToErrorResponse(err.Error())
}
```

We can log the error straight away in the `log.Error(err.Error())` line. The next step is to "dispatch" the incoming error to one of the previously registered error types, in which it got wrapped potentially via `httpError, _ := httpErrorFromError(err)`. Therefore we had to register all error types in the `var httpErrors = make(map[HttpError]bool)` map and can now dispatch via:

```Go
func httpErrorFromError(wrappedError error) (HttpError, bool) {
	for httpError := range httpErrors {
		if errors.Is(wrappedError, httpError) {
			return httpError, true
		}
	}
	return ErrInternalServerError, false
}
```

The idea is to "dispatch" an incoming error to a previously registered custom error or return an `ErrInternalServerErr` as fallback. This allows us to make use of the `httpError.ToErrorResponse(err.Error())` method to generate a HTTP response with proper HTTP status code and error ID via

```Go
res.WriteHeaderAndEntity(LogAndReturnHttpError(err))
```

in the error handler.

## Further Resources

* [Rest API Error Handling in Go (Medium)](https://medium.com/@ozdemir.zynl/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd)
