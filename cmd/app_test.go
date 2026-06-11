package cmd

import "testing"

func TestRunReturnsConfigError(t *testing.T) {
	err := run(t.TempDir() + "/missing.toml")
	if err == nil {
		t.Fatal("run() error = nil, want config error")
	}
}
