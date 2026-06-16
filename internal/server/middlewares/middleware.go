package middlewares

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Validator adapts go-playground validator to Echo's validation interface.
type Validator struct {
	Validator *validator.Validate
}

// Validate checks the request DTO and returns an Echo bad request error on
// validation failure.
func (cv *Validator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
