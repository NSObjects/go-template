package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeConfigWithEnvKeepsFileSourceOverrides(t *testing.T) {
	t.Setenv("GO_TEMPLATE_SYSTEM_PORT", ":9999")

	cfg, err := decodeConfigWithEnv([]byte(`
[system]
port = ":9322"
`), "toml", true)
	if err != nil {
		t.Fatalf("decodeConfigWithEnv() error = %v", err)
	}

	if cfg.System.Port != ":9999" {
		t.Fatalf("System.Port = %q, want env override :9999", cfg.System.Port)
	}
}

func TestDecodeConfigWithEnvAppliesJWTOverride(t *testing.T) {
	t.Setenv("GO_TEMPLATE_JWT_EXPIRE", "7200")
	t.Setenv("GO_TEMPLATE_JWT_ENABLED", "true")

	cfg, err := decodeConfigWithEnv([]byte(`
	[jwt]
	enabled = false
	expire = 3600
	`), "toml", true)
	if err != nil {
		t.Fatalf("decodeConfigWithEnv() error = %v", err)
	}

	if !cfg.JWT.Enabled {
		t.Fatal("JWT.Enabled = false, want env override true")
	}
	if cfg.JWT.Expire != 7200 {
		t.Fatalf("JWT.Expire = %d, want env override 7200", cfg.JWT.Expire)
	}
}

func TestLoadReadsFileAndAppliesEnvOverrides(t *testing.T) {
	t.Setenv("GO_TEMPLATE_SYSTEM_PORT", ":9999")

	path := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(path, []byte(`
[system]
port = ":9322"
level = 1
env = "test"

[jwt]
enabled = false
expire = 3600
`), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.System.Port != ":9999" {
		t.Fatalf("System.Port = %q, want env override :9999", cfg.System.Port)
	}
	if cfg.System.Level != DebugLevel {
		t.Fatalf("System.Level = %d, want %d", cfg.System.Level, DebugLevel)
	}
	if cfg.JWT.Expire != 3600 {
		t.Fatalf("JWT.Expire = %d, want 3600", cfg.JWT.Expire)
	}
}
