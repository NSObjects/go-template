package configs

import "testing"

func TestMergeEnablesConfiguredInfrastructureFlags(t *testing.T) {
	dst := Config{}
	src := Config{
		Mysql:   MysqlConfig{Enabled: true},
		Redis:   RedisConfig{Enabled: true},
		Mongodb: Mongodb{Enabled: true},
		Kafka:   KafkaConfig{Enabled: true},
	}

	got := Merge(dst, src)

	if !got.Mysql.Enabled {
		t.Fatal("Mysql.Enabled = false, want true")
	}
	if !got.Redis.Enabled {
		t.Fatal("Redis.Enabled = false, want true")
	}
	if !got.Mongodb.Enabled {
		t.Fatal("Mongodb.Enabled = false, want true")
	}
	if !got.Kafka.Enabled {
		t.Fatal("Kafka.Enabled = false, want true")
	}
}

func TestMergeDoesNotDisableAlreadyEnabledInfrastructureFlags(t *testing.T) {
	dst := Config{
		Mysql:   MysqlConfig{Enabled: true},
		Redis:   RedisConfig{Enabled: true},
		Mongodb: Mongodb{Enabled: true},
		Kafka:   KafkaConfig{Enabled: true},
	}
	src := Config{}

	got := Merge(dst, src)

	if !got.Mysql.Enabled || !got.Redis.Enabled || !got.Mongodb.Enabled || !got.Kafka.Enabled {
		t.Fatalf("Merge disabled a database capability: %#v", got)
	}
}
