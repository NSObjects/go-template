package configs

import "testing"

func TestDecodeConfigPreservesCapabilityProviderKeys(t *testing.T) {
	cfg, err := decodeConfig([]byte(`
[capabilities.providers]
"user.storage" = "mysql"
`), "toml")
	if err != nil {
		t.Fatalf("decodeConfig() error = %v", err)
	}

	if got := cfg.Capabilities.Providers["user.storage"]; got != "mysql" {
		t.Fatalf(`Capabilities.Providers["user.storage"] = %q, want mysql`, got)
	}
}

func TestDecodeConfigWithEnvKeepsFileSourceOverrides(t *testing.T) {
	t.Setenv("ECHOADMIN_SYSTEM_PORT", ":9999")

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
