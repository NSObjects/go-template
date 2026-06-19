// Package httpreq parses and validates Echo request inputs for HTTP adapters.
package httpreq

import (
	"strconv"

	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/platform/apperr"
)

// BindAndValidate binds a JSON request body and runs the server validator.
func BindAndValidate(c *echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return apperr.WrapBadRequest(err, "invalid request body")
	}
	return c.Validate(req)
}

// PathID parses a positive int64 path parameter.
func PathID(c *echo.Context, name, label string) (int64, error) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		return 0, apperr.NewBadRequest("invalid " + label + " id")
	}
	return id, nil
}

// Pagination parses standard page and page_size query parameters.
func Pagination(c *echo.Context, defaultPageSize int) (int, int, error) {
	page, err := QueryInt(c, "page", 1)
	if err != nil {
		return 0, 0, err
	}
	pageSize, err := QueryInt(c, "page_size", defaultPageSize)
	if err != nil {
		return 0, 0, err
	}
	return page, pageSize, nil
}

// QueryInt parses an optional int query parameter.
func QueryInt(c *echo.Context, name string, fallback int) (int, error) {
	raw := c.QueryParam(name)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, apperr.NewBadRequest("invalid " + name)
	}
	return value, nil
}

// QueryInt64 parses an optional int64 query parameter.
func QueryInt64(c *echo.Context, name string, fallback int64) (int64, error) {
	raw := c.QueryParam(name)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, apperr.NewBadRequest("invalid " + name)
	}
	return value, nil
}

// QueryBool parses an optional bool query parameter.
func QueryBool(c *echo.Context, name string, fallback bool) (bool, error) {
	raw := c.QueryParam(name)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, apperr.NewBadRequest("invalid " + name)
	}
	return value, nil
}
