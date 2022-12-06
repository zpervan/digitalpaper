package core

import (
	"fmt"
	"regexp"
)

// Recommended mail pattern matching (https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address)
var StandardMailAddressPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Tag               string
	ValidationResults []ValidationResult
}

type ValidationResult struct {
	Attribute string `json:"attribute"`
	Error     string `json:"error"`
}

func (v *Validator) MinLength(value *string, minLength int) {
	isValidLength := len(*value) >= minLength

	if !isValidLength {
		errorMessage := fmt.Errorf("value is too short - should contain at least %d characters", minLength)
		v.ValidationResults = append(v.ValidationResults, ValidationResult{v.Tag, errorMessage.Error()})
	}
}

func (v *Validator) MaxLength(value *string, maxLength int) {
	isValidLength := len(*value) <= maxLength

	if !isValidLength {
		errorMessage := fmt.Errorf("value is too long - should contain at least %d characters", maxLength)
		v.ValidationResults = append(v.ValidationResults, ValidationResult{v.Tag, errorMessage.Error()})
	}
}

func (v *Validator) MatchesPattern(value *string, pattern *regexp.Regexp) {
	isValidMail := pattern.MatchString(*value)

	if !isValidMail {
		v.ValidationResults = append(v.ValidationResults, ValidationResult{v.Tag, "value is not valid"})
	}
}

func (v Validator) IsValid() bool {
	return v.ValidationResults == nil
}

func ValidateUser(user *User) []ValidationResult {
	validator := Validator{}

	validator.Tag = "username"
	validator.MinLength(&user.Username, 2)
	validator.MaxLength(&user.Username, 15)

	validator.Tag = "name"
	validator.MinLength(&user.Name, 2)
	validator.MaxLength(&user.Name, 25)

	validator.Tag = "surname"
	validator.MinLength(&user.Surname, 2)
	validator.MaxLength(&user.Surname, 25)

	validator.Tag = "password"
	validator.MinLength(&user.Password, 6)
	validator.MaxLength(&user.Password, 25)

	validator.Tag = "mail"
	validator.MatchesPattern(&user.Mail, StandardMailAddressPattern)

	return validator.ValidationResults
}
