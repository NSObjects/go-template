package configs

// Merge 将 src 非零值覆盖到 dst（浅合并，适合配置）
func Merge(dst, src Config) Config {
	// System
	if src.System.Port != "" {
		dst.System.Port = src.System.Port
	}
	if src.System.Level != 0 {
		dst.System.Level = src.System.Level
	}
	// Log
	if src.Log.Level != "" {
		dst.Log.Level = src.Log.Level
	}
	if src.Log.Format != "" {
		dst.Log.Format = src.Log.Format
	}
	// Console
	if src.Log.Console.Format != "" {
		dst.Log.Console.Format = src.Log.Console.Format
	}
	if src.Log.Console.Output != "" {
		dst.Log.Console.Output = src.Log.Console.Output
	}
	// File
	if src.Log.File.Filename != "" {
		dst.Log.File.Filename = src.Log.File.Filename
	}
	if src.Log.File.MaxSize != 0 {
		dst.Log.File.MaxSize = src.Log.File.MaxSize
	}
	if src.Log.File.MaxBackups != 0 {
		dst.Log.File.MaxBackups = src.Log.File.MaxBackups
	}
	if src.Log.File.MaxAge != 0 {
		dst.Log.File.MaxAge = src.Log.File.MaxAge
	}
	if src.Log.File.Compress {
		dst.Log.File.Compress = true
	}
	if src.Log.File.Format != "" {
		dst.Log.File.Format = src.Log.File.Format
	}
	// Elasticsearch
	if src.Log.Elasticsearch.URL != "" {
		dst.Log.Elasticsearch.URL = src.Log.Elasticsearch.URL
	}
	if src.Log.Elasticsearch.Index != "" {
		dst.Log.Elasticsearch.Index = src.Log.Elasticsearch.Index
	}
	if src.Log.Elasticsearch.Timeout != 0 {
		dst.Log.Elasticsearch.Timeout = src.Log.Elasticsearch.Timeout
	}
	// Loki
	if src.Log.Loki.URL != "" {
		dst.Log.Loki.URL = src.Log.Loki.URL
	}
	if len(src.Log.Loki.Labels) > 0 {
		dst.Log.Loki.Labels = src.Log.Loki.Labels
	}
	if src.Log.Loki.Timeout != 0 {
		dst.Log.Loki.Timeout = src.Log.Loki.Timeout
	}
	// Mysql
	if src.Mysql.Host != "" {
		dst.Mysql.Host = src.Mysql.Host
	}
	if src.Mysql.Port != "" {
		dst.Mysql.Port = src.Mysql.Port
	}
	if src.Mysql.User != "" {
		dst.Mysql.User = src.Mysql.User
	}
	if src.Mysql.Password != "" {
		dst.Mysql.Password = src.Mysql.Password
	}
	if src.Mysql.Database != "" {
		dst.Mysql.Database = src.Mysql.Database
	}
	if src.Mysql.DockerHost != "" {
		dst.Mysql.DockerHost = src.Mysql.DockerHost
	}
	if src.Mysql.MaxOpenConns != 0 {
		dst.Mysql.MaxOpenConns = src.Mysql.MaxOpenConns
	}
	if src.Mysql.MaxIdleConns != 0 {
		dst.Mysql.MaxIdleConns = src.Mysql.MaxIdleConns
	}
	// Redis
	if src.Redis.Host != "" {
		dst.Redis.Host = src.Redis.Host
	}
	if src.Redis.Port != "" {
		dst.Redis.Port = src.Redis.Port
	}
	if src.Redis.Password != "" {
		dst.Redis.Password = src.Redis.Password
	}
	if src.Redis.Database != 0 {
		dst.Redis.Database = src.Redis.Database
	}
	// Mongo
	if src.Mongodb.Host != "" {
		dst.Mongodb.Host = src.Mongodb.Host
	}
	if src.Mongodb.Port != "" {
		dst.Mongodb.Port = src.Mongodb.Port
	}
	if src.Mongodb.User != "" {
		dst.Mongodb.User = src.Mongodb.User
	}
	if src.Mongodb.Password != "" {
		dst.Mongodb.Password = src.Mongodb.Password
	}
	if src.Mongodb.DataBase != "" {
		dst.Mongodb.DataBase = src.Mongodb.DataBase
	}
	// Kafka
	if len(src.Kafka.Brokers) > 0 {
		dst.Kafka.Brokers = src.Kafka.Brokers
	}
	if src.Kafka.ClientID != "" {
		dst.Kafka.ClientID = src.Kafka.ClientID
	}
	if src.Kafka.Topic != "" {
		dst.Kafka.Topic = src.Kafka.Topic
	}
	// JWT
	if src.JWT.Secret != "" {
		dst.JWT.Secret = src.JWT.Secret
	}
	if src.JWT.Expire != 0 {
		dst.JWT.Expire = src.JWT.Expire
	}
	if len(src.JWT.SkipPaths) > 0 {
		dst.JWT.SkipPaths = src.JWT.SkipPaths
	}
	// CORS
	if len(src.CORS.AllowOrigins) > 0 {
		dst.CORS.AllowOrigins = src.CORS.AllowOrigins
	}
	if len(src.CORS.AllowHeaders) > 0 {
		dst.CORS.AllowHeaders = src.CORS.AllowHeaders
	}
	if len(src.CORS.AllowMethods) > 0 {
		dst.CORS.AllowMethods = src.CORS.AllowMethods
	}
	if src.CORS.AllowCredentials {
		dst.CORS.AllowCredentials = true
	}
	// Casbin
	if src.Casbin.Model != "" {
		dst.Casbin.Model = src.Casbin.Model
	}
	if src.Casbin.ModelFile != "" {
		dst.Casbin.ModelFile = src.Casbin.ModelFile
	}
	// Etcd/Consul 客户端配置
	if len(src.Etcd.Endpoints) > 0 {
		dst.Etcd.Endpoints = src.Etcd.Endpoints
	}
	if src.Etcd.Key != "" {
		dst.Etcd.Key = src.Etcd.Key
	}
	if src.Etcd.Format != "" {
		dst.Etcd.Format = src.Etcd.Format
	}
	if src.Etcd.Username != "" {
		dst.Etcd.Username = src.Etcd.Username
	}
	if src.Etcd.Password != "" {
		dst.Etcd.Password = src.Etcd.Password
	}
	if src.Etcd.DialTimeoutSeconds != 0 {
		dst.Etcd.DialTimeoutSeconds = src.Etcd.DialTimeoutSeconds
	}
	if src.Consul.Address != "" {
		dst.Consul.Address = src.Consul.Address
	}
	if src.Consul.Token != "" {
		dst.Consul.Token = src.Consul.Token
	}
	if src.Consul.Key != "" {
		dst.Consul.Key = src.Consul.Key
	}
	if src.Consul.Format != "" {
		dst.Consul.Format = src.Consul.Format
	}
	return dst
}
