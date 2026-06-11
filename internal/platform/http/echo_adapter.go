package http

import (
	"errors"
	"fmt"

	"github.com/NSObjects/go-template/internal/platform/module"
	"github.com/labstack/echo/v4"
)

// ErrUnsupportedEntryPoint identifies a non-HTTP entry point passed to the HTTP adapter.
var ErrUnsupportedEntryPoint = errors.New("unsupported http entry point")

// RoutesFromEntryPoints extracts HTTP routes from generic module entry points.
func RoutesFromEntryPoints(entryPoints []module.EntryPoint) ([]Route, error) {
	routes := make([]Route, 0, len(entryPoints))
	for _, entryPoint := range entryPoints {
		if entryPoint.Type != EntryPointType {
			return nil, fmt.Errorf("%w: type %q for module %q", ErrUnsupportedEntryPoint, entryPoint.Type, entryPoint.Owner)
		}
		route, ok := entryPoint.Value.(Route)
		if !ok {
			return nil, fmt.Errorf("http entry point %q for module %q has invalid route value", entryPoint.Name, entryPoint.Owner)
		}
		if route.Owner == "" {
			route.Owner = entryPoint.Owner
		}
		if route.Name == "" {
			route.Name = entryPoint.Name
		}
		routes = append(routes, route)
	}
	return routes, nil
}

// RegisterRoutes registers HTTP routes on the provided Echo group.
func RegisterRoutes(group *echo.Group, routes []Route) error {
	for _, route := range routes {
		if route.Method == "" {
			return fmt.Errorf("http route %q for module %q has empty method", route.Name, route.Owner)
		}
		if route.Path == "" {
			return fmt.Errorf("http route %q for module %q has empty path", route.Name, route.Owner)
		}
		if route.Handler == nil {
			return fmt.Errorf("http route %q for module %q has nil handler", route.Name, route.Owner)
		}
		group.Add(route.Method, route.Path, route.Handler, route.Middleware...).Name = route.Name
	}
	return nil
}
