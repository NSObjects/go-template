package user

import (
	"net/http"

	"github.com/NSObjects/go-template/internal/api/biz"
	"github.com/NSObjects/go-template/internal/api/service"
	platformhttp "github.com/NSObjects/go-template/internal/platform/http"
	"github.com/NSObjects/go-template/internal/platform/module"
)

// HTTPEntryPoints declares the user HTTP routes as module entry points.
func HTTPEntryPoints(owner string, useCase biz.UserUseCase) []module.EntryPoint {
	controller := service.NewUserController(useCase)
	routes := []platformhttp.Route{
		{
			Owner:   owner,
			Name:    "list users",
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: controller.ListUsers,
		},
		{
			Owner:   owner,
			Name:    "create user",
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: controller.Create,
		},
		{
			Owner:   owner,
			Name:    "get user",
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: controller.GetByID,
		},
		{
			Owner:   owner,
			Name:    "update user",
			Method:  http.MethodPut,
			Path:    "/users/:id",
			Handler: controller.Update,
		},
		{
			Owner:   owner,
			Name:    "delete user",
			Method:  http.MethodDelete,
			Path:    "/users/:id",
			Handler: controller.Delete,
		},
	}

	entryPoints := make([]module.EntryPoint, 0, len(routes))
	for _, route := range routes {
		entryPoints = append(entryPoints, platformhttp.EntryPoint(owner, route))
	}
	return entryPoints
}
