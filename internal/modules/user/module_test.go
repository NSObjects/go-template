package user

import (
	"net/http"
	"testing"

	platformhttp "github.com/NSObjects/go-template/internal/platform/http"
	"github.com/NSObjects/go-template/internal/platform/module"
)

func TestDescriptorDeclaresUserStorageAndHTTPRoutes(t *testing.T) {
	descriptor := New(fakeUserUseCase{}).Descriptor()

	if descriptor.Name != ModuleName {
		t.Fatalf("descriptor.Name = %q, want %q", descriptor.Name, ModuleName)
	}
	if descriptor.Kind != module.BusinessModule {
		t.Fatalf("descriptor.Kind = %q, want business", descriptor.Kind)
	}
	if len(descriptor.Requires) != 1 {
		t.Fatalf("len(descriptor.Requires) = %d, want 1", len(descriptor.Requires))
	}
	if descriptor.Requires[0].Name != StorageCapability {
		t.Fatalf("descriptor.Requires[0].Name = %q, want %q", descriptor.Requires[0].Name, StorageCapability)
	}
	if descriptor.Requires[0].Name == "mysql" || descriptor.Requires[0].Name == "mongodb" {
		t.Fatalf("user module requires concrete storage provider %q", descriptor.Requires[0].Name)
	}

	wantRoutes := map[string]struct {
		method string
		path   string
	}{
		"list users":  {method: http.MethodGet, path: "/users"},
		"create user": {method: http.MethodPost, path: "/users"},
		"get user":    {method: http.MethodGet, path: "/users/:id"},
		"update user": {method: http.MethodPut, path: "/users/:id"},
		"delete user": {method: http.MethodDelete, path: "/users/:id"},
	}
	if len(descriptor.EntryPoints) != len(wantRoutes) {
		t.Fatalf("len(descriptor.EntryPoints) = %d, want %d", len(descriptor.EntryPoints), len(wantRoutes))
	}
	for _, entryPoint := range descriptor.EntryPoints {
		if entryPoint.Owner != ModuleName {
			t.Fatalf("entryPoint.Owner = %q, want %q", entryPoint.Owner, ModuleName)
		}
		if entryPoint.Type != module.EntryPointHTTP {
			t.Fatalf("entryPoint.Type = %q, want http", entryPoint.Type)
		}
		route, ok := entryPoint.Value.(platformhttp.Route)
		if !ok {
			t.Fatalf("entryPoint.Value for %q is not an HTTP route", entryPoint.Name)
		}
		want, ok := wantRoutes[route.Name]
		if !ok {
			t.Fatalf("unexpected route %q", route.Name)
		}
		if route.Method != want.method || route.Path != want.path {
			t.Fatalf("route %q = %s %s, want %s %s", route.Name, route.Method, route.Path, want.method, want.path)
		}
		delete(wantRoutes, route.Name)
	}
	if len(wantRoutes) != 0 {
		t.Fatalf("missing routes: %+v", wantRoutes)
	}
}

var _ UseCase = fakeUserUseCase{}

type fakeUserUseCase struct {
	UseCase
}
