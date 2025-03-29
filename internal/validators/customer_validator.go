package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type CustomerValidator struct {
	validate validator.Validate
}

func NewCustomerValidator() *CustomerValidator {
	v := validator.New()
	_ = v.RegisterValidation("username", validateUsername)
	_ = v.RegisterValidation("password", validatePassword)
	return &CustomerValidator{
		validate: *v,
	}
}

func (cv *CustomerValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) == 0 {
		return false
	}
	letterRegex := regexp.MustCompile(`[a-zA-Z]`)
	hasLetter := letterRegex.MatchString(username)
	numberRegex := regexp.MustCompile(`[0-9]`)
	hasNumber := numberRegex.MatchString(username)

	return hasLetter && hasNumber

}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) == 0 {
		return false
	}
	letterRegex := regexp.MustCompile(`[a-zA-Z]`)
	hasLetter := letterRegex.MatchString(password)

	numberRegex := regexp.MustCompile(`[0-9]`)
	hasNumber := numberRegex.MatchString(password)

	return hasLetter && hasNumber
}
