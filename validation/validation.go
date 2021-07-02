package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/jenpaff/golang-microservices/errors"
	"reflect"
	"strings"
)

const (
	validRegexInput = "validRegexInput"
	notBlank        = "notBlank"
	lte             = "lte"
	gte             = "gte"
	required        = "required"
)

func NewValidate() (*validator.Validate, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// register name defined with the `json` tag, such that we can return the json name with the validation message
		// by default the name of the struct is returned which we do not want to expose to the user
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	err := registerValidation(validate, validRegexInput, isInputValid)
	if err != nil {
		return nil, err
	}
	err = registerValidation(validate, notBlank, validators.NotBlank)
	if err != nil {
		return nil, err
	}
	return validate, nil
}

func registerValidation(validate *validator.Validate, tag string, fn validator.Func) error {
	err := validate.RegisterValidation(tag, fn)
	if err != nil {
		return fmt.Errorf("cannot register tag %s: %w", tag, err)
	}
	return nil
}

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
	case required:
		return fmt.Sprintf("The field %s is required and should not be empty.", field)
	case validRegexInput:
		return fmt.Sprintf("The field %s contains invalid characters- only the following characters are allowed: a-zA-Z0-9", field)
	case lte:
		return fmt.Sprintf("The field %s contains too many characters.", field)
	case gte:
		return fmt.Sprintf("The field %s does not contain enough characters.", field)
	case notBlank:
		return fmt.Sprintf("The field %s should not contain blank spaces only.", field)
	default:
		return fmt.Sprintf("Could not resolve validation tag %s for %s", tag, field)
	}
}
