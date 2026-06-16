package boot

import (
	"context"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
)

func TestRunReturnsConfigLoadError(t *testing.T) {
	err := Run(t.TempDir() + "/missing.toml")
	if err == nil {
		t.Fatal("Run() error = nil, want config load error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("Run() error = %q, want load config context", err.Error())
	}
}

func TestNewAppAssemblesServer(t *testing.T) {
	app, err := NewApp(configs.Config{})
	if err != nil {
		t.Fatalf("NewApp() error = %v", err)
	}
	if app.Server() == nil {
		t.Fatal("Server() = nil, want assembled server")
	}
}

func TestNewAppReturnsConfigError(t *testing.T) {
	app, err := NewApp(configs.Config{
		JWT: configs.JWTConfig{
			Enabled: true,
		},
	})

	if err == nil {
		t.Fatal("NewApp() error = nil, want config error")
	}
	if app != nil {
		t.Fatalf("NewApp() app = %#v, want nil", app)
	}
}

func TestAppRunRejectsNilServer(t *testing.T) {
	err := (&App{}).Run(context.Background())
	if err == nil {
		t.Fatal("Run() error = nil, want nil server error")
	}
}
