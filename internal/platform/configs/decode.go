package configs

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	configFormatJSON = "json"
	configFormatTOML = "toml"
	configFormatYAML = "yaml"
	configFormatYML  = "yml"
)

func decodeConfigWithEnv(data []byte, format string, useEnv bool) (Config, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType(configType(format))
	if useEnv {
		v.SetEnvPrefix(EnvPrefix)
		v.AutomaticEnv()
		v.SetEnvKeyReplacer(strings.NewReplacer("::", "_", ".", "_"))
		if err := bindConfigEnv(v); err != nil {
			return Config{}, err
		}
	}
	if err := v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return Config{}, err
	}
	if useEnv {
		applyEnvListOverrides(&cfg)
	}
	cfg = Normalize(cfg)
	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Normalize applies application config defaults.
func Normalize(cfg Config) Config {
	cfg = normalizeAppDefaults(cfg)
	cfg = normalizeLogDefaults(cfg)
	cfg = normalizeJWTDefaults(cfg)
	cfg = normalizeResourceDefaults(cfg)
	return normalizeTracingDefaults(cfg)
}

func normalizeAppDefaults(cfg Config) Config {
	if cfg.App.Name == "" {
		cfg.App.Name = DefaultAppName
	}
	if cfg.App.Version == "" {
		cfg.App.Version = DefaultAppVersion
	}
	if cfg.System.Port == "" {
		cfg.System.Port = DefaultPort
	}
	if cfg.System.Level == 0 {
		cfg.System.Level = OnlineLevel
	}
	return cfg
}

func normalizeLogDefaults(cfg Config) Config {
	cfg.Log.Format = strings.ToLower(strings.TrimSpace(cfg.Log.Format))
	if cfg.Log.Format == "" {
		cfg.Log.Format = LogFormatJSON
		if cfg.System.Level == DebugLevel {
			cfg.Log.Format = LogFormatConsole
		}
	}
	if cfg.Log.Format == "text" {
		cfg.Log.Format = LogFormatConsole
	}
	cfg.Log.Output = strings.ToLower(strings.TrimSpace(cfg.Log.Output))
	if cfg.Log.Output == "" {
		cfg.Log.Output = LogOutputStdout
	}
	return cfg
}

func normalizeJWTDefaults(cfg Config) Config {
	if cfg.JWT.SkipPaths == nil {
		cfg.JWT.SkipPaths = []string{"/api/health", "/api/info", "/api/ready"}
	}
	return cfg
}

func normalizeResourceDefaults(cfg Config) Config {
	if cfg.MySQL.MaxOpenConns == 0 {
		cfg.MySQL.MaxOpenConns = DefaultMySQLMaxOpenConns
	}
	if cfg.MySQL.MaxIdleConns == 0 {
		cfg.MySQL.MaxIdleConns = DefaultMySQLMaxIdleConns
	}
	if cfg.MySQL.ConnMaxLifetimeSeconds == 0 {
		cfg.MySQL.ConnMaxLifetimeSeconds = DefaultMySQLConnMaxLifetimeSeconds
	}
	if cfg.MySQL.PingTimeoutSeconds == 0 {
		cfg.MySQL.PingTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.Redis.DialTimeoutSeconds == 0 {
		cfg.Redis.DialTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.Redis.ReadTimeoutSeconds == 0 {
		cfg.Redis.ReadTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.Redis.WriteTimeoutSeconds == 0 {
		cfg.Redis.WriteTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.Redis.PingTimeoutSeconds == 0 {
		cfg.Redis.PingTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.MongoDB.ConnectTimeoutSeconds == 0 {
		cfg.MongoDB.ConnectTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	if cfg.MongoDB.PingTimeoutSeconds == 0 {
		cfg.MongoDB.PingTimeoutSeconds = DefaultCapabilityTimeoutSeconds
	}
	return cfg
}

func normalizeTracingDefaults(cfg Config) Config {
	cfg.Tracing.Protocol = strings.ToLower(strings.TrimSpace(cfg.Tracing.Protocol))
	if cfg.Tracing.Protocol == "" {
		cfg.Tracing.Protocol = DefaultTracingProtocol
	}
	if cfg.Tracing.ShutdownTimeoutSeconds == 0 {
		cfg.Tracing.ShutdownTimeoutSeconds = DefaultTracingShutdownTimeoutSeconds
	}
	return cfg
}

// Validate checks cross-field configuration rules that must fail before the
// HTTP runtime starts.
func Validate(cfg Config) error {
	cfg = Normalize(cfg)

	if err := validateAppConfig(cfg.App); err != nil {
		return err
	}
	if err := validateSystemConfig(cfg.System); err != nil {
		return err
	}
	if err := validateLogConfig(cfg.Log); err != nil {
		return err
	}
	if err := validateJWTConfig(cfg.JWT); err != nil {
		return err
	}
	return validateConfiguredResources(cfg)
}

func applyEnvListOverrides(cfg *Config) {
	if cfg == nil {
		return
	}
	applyEnvList(envName("JWT_SKIP_PATHS"), &cfg.JWT.SkipPaths)
	applyEnvList(envName("HTTP_CORS_ALLOW_ORIGINS"), &cfg.HTTP.CORS.AllowOrigins)
	applyEnvList(envName("HTTP_CORS_ALLOW_METHODS"), &cfg.HTTP.CORS.AllowMethods)
	applyEnvList(envName("HTTP_CORS_ALLOW_HEADERS"), &cfg.HTTP.CORS.AllowHeaders)
	applyEnvList(envName("HTTP_CORS_EXPOSE_HEADERS"), &cfg.HTTP.CORS.ExposeHeaders)
}

func envName(name string) string {
	return EnvPrefix + "_" + name
}

func applyEnvList(name string, target *[]string) {
	raw, ok := os.LookupEnv(name)
	if !ok {
		return
	}
	values := splitEnvList(raw)
	*target = values
}

func splitEnvList(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			values = append(values, value)
		}
	}
	return values
}

func validateAppConfig(cfg AppConfig) error {
	if strings.TrimSpace(cfg.Name) == "" {
		return errors.New("app name is required")
	}
	if strings.TrimSpace(cfg.Version) == "" {
		return errors.New("app version is required")
	}
	return nil
}

func validateSystemConfig(cfg SystemConfig) error {
	switch cfg.Level {
	case DebugLevel, OnlineLevel:
		return nil
	default:
		return errors.New("system level must be 1 (debug) or 2 (online)")
	}
}

func validateLogConfig(cfg LogConfig) error {
	switch cfg.Format {
	case LogFormatConsole, LogFormatJSON:
	default:
		return errors.New("log format must be console or json")
	}
	switch cfg.Output {
	case LogOutputStdout, LogOutputStderr:
		return nil
	default:
		return errors.New("log output must be stdout or stderr")
	}
}

func validateJWTConfig(cfg JWTConfig) error {
	if cfg.Enabled && strings.TrimSpace(cfg.Secret) == "" {
		return errors.New("jwt secret is required when jwt is enabled")
	}
	return nil
}

func validateConfiguredResources(cfg Config) error {
	if err := validateHTTPConfig(cfg.HTTP); err != nil {
		return err
	}
	if err := validateMySQLConfig(cfg.MySQL); err != nil {
		return err
	}
	if err := validateRedisConfig(cfg.Redis); err != nil {
		return err
	}
	if err := validateMongoDBConfig(cfg.MongoDB); err != nil {
		return err
	}
	return validateTracingConfig(cfg.Tracing)
}

func validateHTTPConfig(cfg HTTPConfig) error {
	if cfg.CORS.MaxAgeSeconds < 0 {
		return errors.New("http cors max_age_seconds must not be negative")
	}
	if cfg.CORS.Enabled && len(cfg.CORS.AllowOrigins) == 0 {
		return errors.New("http cors allow_origins are required when cors is enabled")
	}
	if cfg.CORS.Enabled && cfg.CORS.AllowCredentials && hasWildcardOrigin(cfg.CORS.AllowOrigins) {
		return errors.New("http cors allow_origins must not include wildcard when allow_credentials is enabled")
	}
	return nil
}

func hasWildcardOrigin(origins []string) bool {
	for _, origin := range origins {
		if strings.TrimSpace(origin) == "*" {
			return true
		}
	}
	return false
}

func validateMySQLConfig(cfg MySQLConfig) error {
	if cfg.MaxOpenConns < 0 {
		return errors.New("mysql max_open_conns must not be negative")
	}
	if cfg.MaxIdleConns < 0 {
		return errors.New("mysql max_idle_conns must not be negative")
	}
	if cfg.ConnMaxLifetimeSeconds < 0 {
		return errors.New("mysql conn_max_lifetime_seconds must not be negative")
	}
	if cfg.PingTimeoutSeconds < 0 {
		return errors.New("mysql ping_timeout_seconds must not be negative")
	}
	if cfg.Enabled && strings.TrimSpace(cfg.DSN) == "" {
		return errors.New("mysql dsn is required when mysql is enabled")
	}
	if cfg.Enabled && cfg.MaxIdleConns > cfg.MaxOpenConns {
		return errors.New("mysql max_idle_conns must not exceed max_open_conns")
	}
	return nil
}

func validateRedisConfig(cfg RedisConfig) error {
	if cfg.DB < 0 {
		return errors.New("redis db must not be negative")
	}
	if cfg.DialTimeoutSeconds < 0 {
		return errors.New("redis dial_timeout_seconds must not be negative")
	}
	if cfg.ReadTimeoutSeconds < 0 {
		return errors.New("redis read_timeout_seconds must not be negative")
	}
	if cfg.WriteTimeoutSeconds < 0 {
		return errors.New("redis write_timeout_seconds must not be negative")
	}
	if cfg.PingTimeoutSeconds < 0 {
		return errors.New("redis ping_timeout_seconds must not be negative")
	}
	if cfg.Enabled && strings.TrimSpace(cfg.Address) == "" {
		return errors.New("redis address is required when redis is enabled")
	}
	return nil
}

func validateMongoDBConfig(cfg MongoDBConfig) error {
	if cfg.ConnectTimeoutSeconds < 0 {
		return errors.New("mongodb connect_timeout_seconds must not be negative")
	}
	if cfg.PingTimeoutSeconds < 0 {
		return errors.New("mongodb ping_timeout_seconds must not be negative")
	}
	if cfg.Enabled && strings.TrimSpace(cfg.URI) == "" {
		return errors.New("mongodb uri is required when mongodb is enabled")
	}
	return nil
}

func validateTracingConfig(cfg TracingConfig) error {
	if cfg.SampleRatio < 0 || cfg.SampleRatio > 1 {
		return errors.New("tracing sample_ratio must be between 0 and 1")
	}
	if cfg.ShutdownTimeoutSeconds < 0 {
		return errors.New("tracing shutdown_timeout_seconds must not be negative")
	}
	switch cfg.Protocol {
	case "grpc", "http":
	default:
		return errors.New("tracing protocol must be grpc or http")
	}
	if cfg.Enabled && strings.TrimSpace(cfg.Endpoint) == "" {
		return errors.New("tracing endpoint is required when tracing is enabled")
	}
	return nil
}

func configType(format string) string {
	switch strings.ToLower(format) {
	case configFormatJSON:
		return configFormatJSON
	case configFormatYAML, configFormatYML:
		return configFormatYAML
	default:
		return configFormatTOML
	}
}

func bindConfigEnv(v *viper.Viper) error {
	for _, key := range configEnvKeys() {
		if err := v.BindEnv(key); err != nil {
			return err
		}
	}
	return nil
}

func configEnvKeys() []string {
	keys := append([]string{}, coreEnvKeys()...)
	keys = append(keys, httpEnvKeys()...)
	keys = append(keys, resourceEnvKeys()...)
	keys = append(keys, tracingEnvKeys()...)
	return keys
}

func coreEnvKeys() []string {
	return []string{
		"app::name",
		"app::version",
		"system::port",
		"system::level",
		"log::format",
		"log::output",
		"log::caller",
		"jwt::enabled",
		"jwt::secret",
		"jwt::skip_paths",
	}
}

func httpEnvKeys() []string {
	return []string{
		"http::recovery_disabled",
		"http::request_context_disabled",
		"http::request_log_disabled",
		"http::gzip_disabled",
		"http::cors::enabled",
		"http::cors::allow_origins",
		"http::cors::allow_methods",
		"http::cors::allow_headers",
		"http::cors::allow_credentials",
		"http::cors::expose_headers",
		"http::cors::max_age_seconds",
	}
}

func resourceEnvKeys() []string {
	return []string{
		"mysql::enabled",
		"mysql::dsn",
		"mysql::max_open_conns",
		"mysql::max_idle_conns",
		"mysql::conn_max_lifetime_seconds",
		"mysql::ping_timeout_seconds",
		"redis::enabled",
		"redis::address",
		"redis::username",
		"redis::password",
		"redis::db",
		"redis::dial_timeout_seconds",
		"redis::read_timeout_seconds",
		"redis::write_timeout_seconds",
		"redis::ping_timeout_seconds",
		"mongodb::enabled",
		"mongodb::uri",
		"mongodb::database",
		"mongodb::connect_timeout_seconds",
		"mongodb::ping_timeout_seconds",
	}
}

func tracingEnvKeys() []string {
	return []string{
		"tracing::enabled",
		"tracing::endpoint",
		"tracing::protocol",
		"tracing::insecure",
		"tracing::service_name",
		"tracing::sample_ratio",
		"tracing::shutdown_timeout_seconds",
	}
}
