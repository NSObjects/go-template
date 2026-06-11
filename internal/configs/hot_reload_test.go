package configs

import (
	"testing"
	"time"
)

func TestStoreRegisterConfigReloadCallbackReceivesUpdatedConfig(t *testing.T) {
	initial := Config{}
	updated := Config{}
	updated.System.Port = "9090"

	store := NewStore(initial)
	got := make(chan Config, 1)
	store.RegisterConfigReloadCallback(func(config *Config) {
		got <- *config
	})

	store.Update(updated)

	select {
	case config := <-got:
		if config.System.Port != updated.System.Port {
			t.Fatalf("callback config System.Port = %q, want %q", config.System.Port, updated.System.Port)
		}
	case <-time.After(time.Second):
		t.Fatal("callback was not called")
	}
}
