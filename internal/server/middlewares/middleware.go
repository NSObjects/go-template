package middlewares

import (
	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/go-playground/validator/v10"
)

// Validator adapts go-playground validator to Echo's validation interface.
type Validator struct {
	Validator *validator.Validate
}

// Validate checks the request DTO and returns a stable public error while
// keeping validator details in the wrapped cause for logs.
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return apperr.WrapBadRequest(err, "invalid request")
	}
	return nil
}
