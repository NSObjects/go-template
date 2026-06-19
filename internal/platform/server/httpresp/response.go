// Package httpresp renders framework-specific HTTP error responses.
package httpresp

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"

	"github.com/NSObjects/go-template/internal/platform/apperr"
	"github.com/NSObjects/go-template/internal/platform/requestctx"
)

// ErrorResponse is the standard JSON error response for HTTP adapters.
type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
	Timestamp int64  `json:"timestamp"`
}

// Response is the standard JSON success response for HTTP adapters.
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id"`
	Timestamp int64       `json:"timestamp"`
}

// PageMeta describes paginated list response metadata.
type PageMeta struct {
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	HasNext  bool `json:"has_next"`
}

// ListResponse is the standard JSON list response for HTTP adapters.
type ListResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Page      PageMeta    `json:"page"`
	RequestID string      `json:"request_id"`
	Timestamp int64       `json:"timestamp"`
}

// OK renders a successful response envelope.
func OK(c *echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:      apperr.ErrSuccess,
		Message:   "OK",
		Data:      data,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	})
}

// Created renders a successful creation response envelope.
func Created(c *echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Code:      apperr.ErrSuccess,
		Message:   "OK",
		Data:      data,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	})
}

// List renders a successful paginated list response envelope.
func List(c *echo.Context, data interface{}, page PageMeta) error {
	return c.JSON(http.StatusOK, ListResponse{
		Code:      apperr.ErrSuccess,
		Message:   "OK",
		Data:      data,
		Page:      page,
		RequestID: RequestID(c),
		Timestamp: time.Now().Unix(),
	})
}

// NewPageMeta validates and creates pagination metadata.
func NewPageMeta(page, pageSize, total int) (PageMeta, error) {
	if page < 1 {
		return PageMeta{}, apperr.NewBadRequest("invalid pagination")
	}
	if pageSize < 1 || pageSize > 1000 {
		return PageMeta{}, apperr.NewBadRequest("invalid pagination")
	}
	if total < 0 {
		return PageMeta{}, apperr.NewBadRequest("invalid pagination")
	}
	return PageMeta{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		HasNext:  page*pageSize < total,
	}, nil
}

// APIError renders a project error as a JSON HTTP response.
func APIError(c *echo.Context, err error) error {
	if err == nil {
		return errors.New("error cannot be nil")
	}
	if response, unwrapErr := echo.UnwrapResponse(c.Response()); unwrapErr == nil && response.Committed {
		return nil
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
	case apperr.KindMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case apperr.KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// RequestID returns the request ID used in HTTP error responses.
func RequestID(c *echo.Context) string {
	if requestID := requestctx.GetRequestID(c.Request().Context()); requestID != "" {
		setResponseRequestID(c, requestID)
		return requestID
	}

	if requestID := requestctx.CleanMetadataID(c.Request().Header.Get("X-Request-ID")); requestID != "" {
		setResponseRequestID(c, requestID)
		return requestID
	}

	if requestID := requestctx.CleanMetadataID(c.Response().Header().Get("X-Request-ID")); requestID != "" {
		return requestID
	}

	requestID := generateRequestID()
	setResponseRequestID(c, requestID)
	return requestID
}

func setResponseRequestID(c *echo.Context, requestID string) {
	c.Response().Header().Set("X-Request-ID", requestID)
}

func generateRequestID() string {
	return uuid.NewString()
}
