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
	UserName    string `json:"name" validate:"required,validRegexInput"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
```

In this example we're using 2 validation tags: 
1. `required`: defines the field as required. This field is already defined by the `go-validate` package
2. `validRegexInput`: this is a custom validation tag, implemented by us, so we can validate the input against a certain regex. 

To make sure our struct is validated and meaningful validation errors returned, we add this snippet to our `user-controller.go`

```go
err = c.validator.Struct(creationRequest)
if err != nil {
    return validation.GetValidationError(err)
}
```

## Code

### Implementation
The code can be found in the `/validation` package. 

A tag such as `required` doesn't need to be registered, however, our custom tag `invalidRegexInput` 
needs to be registered such that our validator can recognise the custom tag. 

This happens when we initialise our validator: 

```go
func NewValidate() (*validator.Validate, error) {
	validate := validator.New()
	err := registerValidation(validate, "validRegexInput", isInputValid)
	if err != nil {
		return nil, err
	}
	return validate, nil
}
```

This registration will map the tag to a func we have defined in `inputregex.go`, specifically we only allow
alphanumeric fields as valid input:

```go
var inputRegex = regexp.MustCompile(`^[\w]*$`)

func isInputValid(fl validator.FieldLevel) bool {
	field := fl.Field()
	fieldType := field.Kind()
	if fieldType == reflect.String {
		return inputRegex.MatchString(field.String())
	}
	log.Errorf("Could not match type of field %s ", fieldType)
	return false
}
```

### Improvements

**Problem**: When returning a validation error the library would specify the struct name e.g. 
`Validation error in UserName` rather than `Validation error in name` when looking at our `UserCreationRequest` struct.

**Solution**: Add the following snippet to our validation initialisation

```go
validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
    // register name defined with the `json` tag, such that we can return the json name with the validation message
    // by default the name of the struct is returned which we do not want to expose to the user
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if name == "-" {
        return ""
    }
    return name
})
```


**Problem**: On testing our implementation, we found that the default messages given by the library weren't concise enough for our users.

**Solution**: We implemented custom messages, as follows: 

```go
func GetValidationError(err error) error {
	if e, ok := err.(validator.ValidationErrors); ok {
		var errorMessages = make([]string, len(e))
		for i, fieldError := range e {
			errorMessages[i] = getMessage(fieldError)
		}
		error := strings.Join(errorMessages, ";")
		return fmt.Errorf("Error when validating request body: %s: %w", error, errors.InvalidInput)
	} else {
		return fmt.Errorf("Error validating request %s : %w", err.Error(), errors.InvalidInput)
	}
}

func getMessage(fieldError validator.FieldError) string {
	field := fieldError.Field()
	tag := fieldError.Tag()
	switch tag {
	case validRegexInput:
		return fmt.Sprintf("The field %s contains invalid characters- only the following characters are allowed: a-zA-Z0-9", field)
	case required:
		return fmt.Sprintf("The field %s is required.", field)
	default:
		return fmt.Sprintf("Could not resolve validation tag %s for %s", tag, field)
	}
}
```

## Further Resources
- [go-validate package documentation](https://pkg.go.dev/github.com/gookit/validate)