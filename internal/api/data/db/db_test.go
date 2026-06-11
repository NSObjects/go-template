package db

import (
	"context"
	"strings"
	"testing"

	"github.com/NSObjects/go-template/internal/configs"
	"github.com/samber/do/v2"
)

func TestNewDataManagerDoesNotInitializeDisabledComponents(t *testing.T) {
	cfg := configs.Config{
		Mysql: configs.MysqlConfig{
			Host:     "127.0.0.1",
			Port:     "3306",
			User:     "root",
			Password: "password",
			Database: "template",
		},
		Redis: configs.RedisConfig{
			Host: "127.0.0.1",
			Port: "6379",
		},
		Mongodb: configs.Mongodb{
			Host:     "127.0.0.1",
			Port:     "27017",
			DataBase: "template",
		},
		Kafka: configs.KafkaConfig{
			Brokers:  []string{"127.0.0.1:9092"},
			ClientID: "template",
			Topic:    "events",
		},
	}

	dm, err := NewDataManager(cfg)
	if err != nil {
		t.Fatalf("NewDataManager() error = %v", err)
	}

	if dm.Mysql != nil {
		t.Fatal("Mysql was initialized while disabled")
	}
	if dm.Query != nil {
		t.Fatal("Query was initialized while mysql disabled")
	}
	if dm.Redis != nil {
		t.Fatal("Redis was initialized while disabled")
	}
	if dm.Mongodb != nil {
		t.Fatal("Mongodb was initialized while disabled")
	}
	if dm.Kafka != nil {
		t.Fatal("Kafka was initialized while disabled")
	}

	for _, component := range []string{"mysql", "redis", "mongodb", "kafka"} {
		if dm.IsComponentEnabled(component) {
			t.Fatalf("%s is enabled, want disabled", component)
		}
	}

	health := dm.Health(context.Background())
	if len(health) != 0 {
		t.Fatalf("Health() = %#v, want empty health map for disabled components", health)
	}
}

func TestNewDBAndNewQueryReturnNilWhenMysqlDisabled(t *testing.T) {
	dm, err := NewDataManager(configs.Config{})
	if err != nil {
		t.Fatalf("NewDataManager() error = %v", err)
	}

	if got := NewDB(dm); got != nil {
		t.Fatal("NewDB() returned db while mysql disabled")
	}
	if got := NewQuery(dm); got != nil {
		t.Fatal("NewQuery() returned query while mysql disabled")
	}
}

func TestModelStartsWithDisabledComponents(t *testing.T) {
	i := do.New(Register)
	do.ProvideValue(i, configs.Config{
		Mysql: configs.MysqlConfig{
			Host:     "127.0.0.1",
			Port:     "3306",
			User:     "root",
			Password: "password",
			Database: "template",
		},
		Redis: configs.RedisConfig{
			Host: "127.0.0.1",
			Port: "6379",
		},
		Mongodb: configs.Mongodb{
			Host:     "127.0.0.1",
			Port:     "27017",
			DataBase: "template",
		},
		Kafka: configs.KafkaConfig{
			Brokers: []string{"127.0.0.1:9092"},
		},
	})

	dm, err := do.Invoke[*DataManager](i)
	if err != nil {
		t.Fatalf("do.Invoke[*DataManager]() error = %v", err)
	}
	t.Cleanup(func() {
		report := i.Shutdown()
		if !report.Succeed {
			t.Fatalf("injector shutdown failed: %v", report.Errors)
		}
	})

	if dm == nil {
		t.Fatal("DataManager was not populated")
	}
	if len(dm.Health(context.Background())) != 0 {
		t.Fatal("DataManager reported disabled components in health check")
	}
}

func TestNewDataManagerReturnsConfigErrorForEnabledComponents(t *testing.T) {
	tests := []struct {
		name    string
		cfg     configs.Config
		wantErr string
	}{
		{
			name: "mysql enabled with missing fields",
			cfg: configs.Config{
				Mysql: configs.MysqlConfig{Enabled: true},
			},
			wantErr: "mysql enabled",
		},
		{
			name: "redis enabled with missing fields",
			cfg: configs.Config{
				Redis: configs.RedisConfig{Enabled: true},
			},
			wantErr: "redis enabled",
		},
		{
			name: "mongodb enabled with missing fields",
			cfg: configs.Config{
				Mongodb: configs.Mongodb{Enabled: true},
			},
			wantErr: "mongodb enabled",
		},
		{
			name: "kafka enabled with missing brokers",
			cfg: configs.Config{
				Kafka: configs.KafkaConfig{Enabled: true},
			},
			wantErr: "kafka enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm, err := NewDataManager(tt.cfg)
			if err == nil {
				t.Fatal("NewDataManager() error = nil, want config error")
			}
			if dm != nil {
				t.Fatal("NewDataManager() returned DataManager with invalid config")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("NewDataManager() error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
