package middlewares

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/NSObjects/go-template/internal/platform/apperr"
)

// Validator adapts go-playground validator to Echo's validation interface.
type Validator struct {
	Validator *validator.Validate
}

// Validate checks the request DTO and returns a stable public error while
// keeping validator details in the wrapped cause for logs.
func (cv *Validator) Validate(i interface{}) error {
	if cv == nil || cv.Validator == nil {
		return apperr.WrapInternal(errors.New("validator is not configured"), "internal server error")
	}
	if err := cv.Validator.Struct(i); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
			fieldError := validationErrors[0]
			return apperr.NewValidation(fieldError.Field(), validationMessage(fieldError))
		}
		return apperr.WrapBadRequest(err, "invalid request")
	}
	return nil
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	default:
		return fmt.Sprintf("%s failed %s validation", err.Field(), err.Tag())
	}
}
