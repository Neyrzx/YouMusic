package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

const (
	serviceKey  = "service"
	serviceName = "youmusic"
)

//nolint:gochecknoglobals //исключение для логгера
var (
	defaultLogger     *zerolog.Logger
	defaultLoggerOnce sync.Once
)

func NewLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Str(serviceKey, serviceName).
		Timestamp().
		Caller().
		Logger()

	return &logger
}

func DefaultLogger() *zerolog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}
