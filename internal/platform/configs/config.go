// Package configs loads and validates static application configuration.
package configs

// Level identifies the runtime logging and debug mode.
type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota + 1
	// OnlineLevel is the default production priority.
	OnlineLevel
)

const (
	// DefaultAppName is returned by /api/info when config omits app.name.
	DefaultAppName = "go-template"

	// EnvPrefix is the environment variable prefix used for config overrides.
	EnvPrefix = "GO_TEMPLATE"

	// DefaultAppVersion is returned by /api/info when config omits app.version.
	DefaultAppVersion = "dev"

	// DefaultPort is the HTTP port used when config omits system.port.
	DefaultPort = ":9322"

	// LogFormatConsole writes human-readable zerolog console output.
	LogFormatConsole = "console"

	// LogFormatJSON writes structured zerolog JSON output.
	LogFormatJSON = "json"

	// LogOutputStdout writes logs to standard output.
	LogOutputStdout = "stdout"

	// LogOutputStderr writes logs to standard error.
	LogOutputStderr = "stderr"

	// DefaultCapabilityTimeoutSeconds is the default health-check timeout for
	// external resources.
	DefaultCapabilityTimeoutSeconds = 3

	// DefaultMySQLMaxOpenConns is the default upper bound for the MySQL pool.
	DefaultMySQLMaxOpenConns = 25

	// DefaultMySQLMaxIdleConns is the default idle connection count for MySQL.
	DefaultMySQLMaxIdleConns = 5

	// DefaultMySQLConnMaxLifetimeSeconds controls how long MySQL connections can be reused.
	DefaultMySQLConnMaxLifetimeSeconds = 300

	// DefaultRedisDB is the default logical Redis database.
	DefaultRedisDB = 0

	// DefaultTracingProtocol is the OTLP transport used by default.
	DefaultTracingProtocol = "grpc"

	// DefaultTracingShutdownTimeoutSeconds bounds provider flush during shutdown.
	DefaultTracingShutdownTimeoutSeconds = 5
)

// Config is the complete application configuration loaded at startup.
type Config struct {
	App     AppConfig     `mapstructure:"app"`
	System  SystemConfig  `mapstructure:"system"`
	Log     LogConfig     `mapstructure:"log"`
	HTTP    HTTPConfig    `mapstructure:"http"`
	JWT     JWTConfig     `mapstructure:"jwt"`
	MySQL   MySQLConfig   `mapstructure:"mysql"`
	Redis   RedisConfig   `mapstructure:"redis"`
	MongoDB MongoDBConfig `mapstructure:"mongodb"`
	Tracing TracingConfig `mapstructure:"tracing"`
}

// AppConfig controls process identity exposed by system routes.
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// SystemConfig controls process-level runtime settings.
type SystemConfig struct {
	Port  string `mapstructure:"port"`
	Level Level  `mapstructure:"level"`
}

// LogConfig controls process-level structured logging.
type LogConfig struct {
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Caller bool   `mapstructure:"caller"`
}

// HTTPConfig controls server-owned HTTP middleware.
type HTTPConfig struct {
	RecoveryDisabled       bool       `mapstructure:"recovery_disabled"`
	RequestContextDisabled bool       `mapstructure:"request_context_disabled"`
	RequestLogDisabled     bool       `mapstructure:"request_log_disabled"`
	GzipDisabled           bool       `mapstructure:"gzip_disabled"`
	CORS                   CORSConfig `mapstructure:"cors"`
}

// CORSConfig controls optional CORS middleware.
type CORSConfig struct {
	Enabled          bool     `mapstructure:"enabled"`
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	MaxAgeSeconds    int      `mapstructure:"max_age_seconds"`
}

// JWTConfig controls optional server-level JWT verification.
type JWTConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Secret    string   `mapstructure:"secret"`
	SkipPaths []string `mapstructure:"skip_paths"`
}

// MySQLConfig controls the optional process-level MySQL resource.
type MySQLConfig struct {
	Enabled                bool   `mapstructure:"enabled"`
	DSN                    string `mapstructure:"dsn"`
	MaxOpenConns           int    `mapstructure:"max_open_conns"`
	MaxIdleConns           int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `mapstructure:"conn_max_lifetime_seconds"`
	PingTimeoutSeconds     int    `mapstructure:"ping_timeout_seconds"`
}

// RedisConfig controls the optional process-level Redis resource.
type RedisConfig struct {
	Enabled             bool   `mapstructure:"enabled"`
	Address             string `mapstructure:"address"`
	Username            string `mapstructure:"username"`
	Password            string `mapstructure:"password"`
	DB                  int    `mapstructure:"db"`
	DialTimeoutSeconds  int    `mapstructure:"dial_timeout_seconds"`
	ReadTimeoutSeconds  int    `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `mapstructure:"write_timeout_seconds"`
	PingTimeoutSeconds  int    `mapstructure:"ping_timeout_seconds"`
}

// MongoDBConfig controls the optional process-level MongoDB resource.
type MongoDBConfig struct {
	Enabled               bool   `mapstructure:"enabled"`
	URI                   string `mapstructure:"uri"`
	Database              string `mapstructure:"database"`
	ConnectTimeoutSeconds int    `mapstructure:"connect_timeout_seconds"`
	PingTimeoutSeconds    int    `mapstructure:"ping_timeout_seconds"`
}

// TracingConfig controls optional OpenTelemetry export to Jaeger OTLP.
type TracingConfig struct {
	Enabled                bool    `mapstructure:"enabled"`
	Endpoint               string  `mapstructure:"endpoint"`
	Protocol               string  `mapstructure:"protocol"`
	Insecure               bool    `mapstructure:"insecure"`
	ServiceName            string  `mapstructure:"service_name"`
	SampleRatio            float64 `mapstructure:"sample_ratio"`
	ShutdownTimeoutSeconds int     `mapstructure:"shutdown_timeout_seconds"`
}
