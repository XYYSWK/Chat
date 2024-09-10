package common

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidatorEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	ok, _ := regexp.MatchString(`^\w{5,}@[a-z0-9]{2,3}\.[a-z]+$|\,$`, email)
	return ok
}
