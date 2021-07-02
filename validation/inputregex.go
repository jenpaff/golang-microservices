package validation

import (
	"github.com/go-playground/log"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

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
