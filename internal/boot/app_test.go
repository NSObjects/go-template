package boot

import (
	"strings"
	"testing"
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
