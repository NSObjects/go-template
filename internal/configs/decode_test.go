package configs

import "testing"

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

func TestDecodeConfigWithEnvAppliesMysqlOverride(t *testing.T) {
	t.Setenv("ECHOADMIN_MYSQL_ENABLED", "true")

	cfg, err := decodeConfigWithEnv([]byte(`
	[mysql]
	enabled = false
	`), "toml", true)
	if err != nil {
		t.Fatalf("decodeConfigWithEnv() error = %v", err)
	}

	if !cfg.Mysql.Enabled {
		t.Fatal("Mysql.Enabled = false, want env override true")
	}
}
