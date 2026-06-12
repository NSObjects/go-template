package configs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSourceLoadsCapabilityProviderSelections(t *testing.T) {
	path := writeTempConfig(t, `
[capabilities.providers]
"database.gorm" = "mysql"
`)

	cfg, err := FileSource{Path: path}.Load(context.Background())
	if err != nil {
		t.Fatalf("FileSource.Load() error = %v", err)
	}

	got := cfg.Capabilities.Providers["database.gorm"]
	if got != "mysql" {
		t.Fatalf(`Capabilities.Providers["database.gorm"] = %q, want mysql`, got)
	}
}

func TestFileSourceAllowsOmittedCapabilityProviderSelections(t *testing.T) {
	path := writeTempConfig(t, `
[system]
port = ":9322"
`)

	cfg, err := FileSource{Path: path}.Load(context.Background())
	if err != nil {
		t.Fatalf("FileSource.Load() error = %v", err)
	}

	if len(cfg.Capabilities.Providers) != 0 {
		t.Fatalf("len(Capabilities.Providers) = %d, want 0", len(cfg.Capabilities.Providers))
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config fixture: %v", err)
	}
	return path
}
