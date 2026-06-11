// Package http adapts module-declared HTTP entry points to the active server.
package http

import (
	"github.com/NSObjects/go-template/internal/platform/module"
	"github.com/labstack/echo/v4"
)

const EntryPointType = module.EntryPointHTTP

// Route describes one HTTP entry point exposed by a business module.
type Route struct {
	Owner      string
	Name       string
	Method     string
	Path       string
	Handler    echo.HandlerFunc
	Middleware []echo.MiddlewareFunc
}

// EntryPoint wraps a route as a generic module entry point.
func EntryPoint(owner string, route Route) module.EntryPoint {
	if route.Owner == "" {
		route.Owner = owner
	}
	return module.EntryPoint{
		Owner: route.Owner,
		Type:  EntryPointType,
		Name:  route.Name,
		Value: route,
	}
}
