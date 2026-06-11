package utils

import (
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestNewCasbinEnforcerUsesCasbinV3Adapter(t *testing.T) {
	t.Chdir("../..")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	enforcer, err := NewCasbinEnforcer(db, configs.Config{})
	if err != nil {
		t.Fatalf("NewCasbinEnforcer() error = %v", err)
	}
	if enforcer == nil {
		t.Fatal("NewCasbinEnforcer() returned nil enforcer")
	}
}
