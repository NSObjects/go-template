// Package httpresp renders framework-specific HTTP error responses.
package httpresp

import (
	"errors"
	"net/http"
	"time"

	"github.com/NSObjects/go-template/internal/apperr"
	"github.com/NSObjects/go-template/internal/requestctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ErrorResponse is the standard JSON error response for HTTP adapters.
type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Timestamp int64  `json:"timestamp"`
}

// APIError renders a project error as a JSON HTTP response.
func APIError(c echo.Context, err error) error {
	if err == nil {
		return errors.New("error cannot be nil")
	}

	info := apperr.NewInfo(err)
	rjson := ErrorResponse{
		Code:      info.Code,
		Message:   info.Message,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	}

	return c.JSON(Status(info.Kind), rjson)
}

// Status maps framework-free application error kinds to HTTP status codes.
func Status(kind apperr.Kind) int {
	switch kind {
	case apperr.KindOK:
		return http.StatusOK
	case apperr.KindBadRequest, apperr.KindValidation:
		return http.StatusBadRequest
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// RequestID returns the request ID used in HTTP error responses.
func RequestID(c echo.Context) string {
	if requestID := requestctx.GetRequestID(c.Request().Context()); requestID != "" {
		return requestID
	}

	if requestID := c.Request().Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
		return requestID
	}

	requestID := generateRequestID()
	c.Response().Header().Set("X-Request-ID", requestID)
	return requestID
}

func generateRequestID() string {
	return uuid.NewString()
}
