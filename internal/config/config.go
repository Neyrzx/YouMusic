package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type RunningMode string

const (
	ModeLocal       RunningMode = "local"
	ModeDevelopment RunningMode = "developent"
	ModeProduction  RunningMode = "production"
)

func GetCurrentRunningMode() RunningMode {
	switch RunningMode(os.Getenv("ENVIRONMENT")) {
	case ModeLocal:
		return ModeLocal
	case ModeDevelopment:
		return ModeDevelopment
	case ModeProduction:
		return ModeProduction
	default:
		return ModeLocal
	}
}

type App struct {
	Environment      string `env:"ENVIRONMENT"`
	MigrationsSource string `env:"MIGRATIONS_SOURCE"`
	SwaggerDocPath   string `env:"SWAGGER_DOC_PATH"`
	Server           Server
	GatewayMusicInfo GatewayHTTPClient
	Database         Database
}

type Server struct {
	ServerAddr               string        `env:"SERVER_ADDR"`
	GracefulShoutdownTimeout time.Duration `env:"GRACEFUL_SHUTDOWN_TIMEOUT"`
}

type GatewayHTTPClient struct {
	RetryStratagyDelay       time.Duration `env:"GATEWAY_RETRY_STRATAGY_DELAY"`
	RetryStrategyMaxDelay    time.Duration `env:"GATEWAY_RETRY_STRATAGY_MAX_DELAY"`
	RetryStrategyMaxDuration time.Duration `env:"GATEWAY_RETRY_STRATAGY_MAX_DURATION"`
	RetryStrategyFactor      float64       `env:"GATEWAY_RETRY_STRATAGY_FACTOR"`
	BaseURL                  string        `env:"GATEWAY_MUSIC_INFO_BASE_URL"`
}

// Database представляет собой конфигурацию соединений с базой данных, основанную на переменных окружения.
type Database struct {
	Name              string `env:"DB_NAME, required"`
	Hostname          string `env:"DB_HOST, default=localhost"`
	Port              string `env:"DB_PORT, default=5432"`
	User              string `env:"DB_USER"`
	Password          string `env:"DB_PASSWORD"`
	SSLMode           string `env:"DB_SSLMODE, default=disable"`
	ConnectionTimeout int    `env:"DB_CONNECT_TIMEOUT, default=10"`

	// Proto определяет протокол соединения с базой данных.
	Proto string `env:"DB_PROTO, default=postgresql"`
}

// ConnectionURI формирует и возвращает строку Connection URI на основе конфигурации базы данных.
func (cfg Database) ConnectionURI() string {
	u := url.URL{
		Scheme:   cfg.Proto,
		Host:     cfg.Host(),
		Path:     cfg.Name,
		User:     cfg.UserInfo(),
		RawQuery: cfg.parameterList(),
	}

	// TODO: валидацию uri

	return u.String()
}

// Host возвращает строку, содержащую хост и порт для соединения.
func (cfg Database) Host() string {
	if cfg.Port == "" {
		return cfg.Hostname
	}

	return fmt.Sprintf("%s:%s", cfg.Hostname, cfg.Port)
}

// UserInfo возвращает информацию о пользователе для подключения к базе данных.
func (cfg Database) UserInfo() *url.Userinfo {
	if cfg.User != "" && cfg.Password == "" {
		return url.User(cfg.User)
	}

	if cfg.User != "" && cfg.Password != "" {
		return url.UserPassword(cfg.User, cfg.Password)
	}

	return nil
}

// parameterList собирает и возвращает параметры подключения к базе данных.
func (cfg Database) parameterList() string {
	v := make(url.Values)

	if cfg.isValidSSLMode(cfg.SSLMode) {
		v.Add("sslmode", cfg.SSLMode)
	}

	if cfg.isValidConnectTimeout(cfg.ConnectionTimeout) {
		v.Add("connect_timeout", strconv.Itoa(cfg.ConnectionTimeout))
	}

	return v.Encode()
}

// isValidSSLMode проверяет валидность SSLMode.
func (cfg Database) isValidSSLMode(mode string) bool {
	switch mode {
	case "disable", "allow", "prefer", "require", "varify-ca", "verify-full":
		return true
	default:
		return false
	}
}

// isValidConnectTimeout проверяет валидность значения ConnectTimeout.
func (cfg Database) isValidConnectTimeout(value int) bool {
	return value > 0
}
