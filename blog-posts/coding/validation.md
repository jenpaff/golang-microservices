# Validation

## Motivation
- validate format to avoid errors early rather than failing on a db level
- validate any input we receive via POST requests to avoid SQL injection
- having tested a couple of libraries we switched to go-validate when we were forced to add custom regex pattern

## Usage

For an extensive description on how to use the package, please refer to the [go validate package doc](https://pkg.go.dev/github.com/gookit/validate).

Looking at our `POST /users` request, we want to validate the request body before we process any received request. 
To define our input validation we write our input struct like this:

```go
type UserCreationRequest struct {
	UserName    string `json:"name" validate:"notBlank,validRegexInput"`
	Email       string `json:"email" validate:"notBlank"`
	PhoneNumber string `json:"phone_number" validate:"notBlank"`
}
```

In this example we're using 2 validation tags: 
1. `notBlank`: defines the field as required and does not allow blank spaces. This field is already defined by the `go-validate` package
2. `validRegexInput`: this is a custom validation tag, implemented by us, so we can validate the input against a certain regex. 

To make sure our struct is validated and meaningful validation errors returned, we add this snippet to our `user-controller.go`

```go
err = c.validator.Struct(creationRequest)
if err != nil {
    return validation.GetValidationError(err)
}
```

## Code
- [ ] add explanation to the code implementation

## Further Resources
