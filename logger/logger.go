package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level Level
}

type Level string

const (
	TRACE Level = "trace"
	DEBUG Level = "debug"
	INFO  Level = "info"
	WARN  Level = "warn"
	ERROR Level = "error"
	PANIC Level = "panic"
)

func NewLogger(config LoggerConfig) zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: time.RFC3339,
		FormatTimestamp: func(i interface{}) string {
			return fmt.Sprintf("[%s]", i)
		},
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf("[%s]", strings.ToUpper(i.(string)))
		},
		FormatCaller: func(i interface{}) string {
			return filepath.Base(i.(string))
		},
		FormatMessage: func(i interface{}) string {
			return i.(string)
		},
	}).With().Timestamp().Caller().Logger().Level(zerologLevelFromString(config.Level))
}

func zerologLevelFromString(l Level) zerolog.Level {
	switch l {
	case TRACE:
		return zerolog.TraceLevel
	case DEBUG:
		return zerolog.DebugLevel
	case INFO:
		return zerolog.InfoLevel
	case WARN:
		return zerolog.WarnLevel
	case ERROR:
		return zerolog.ErrorLevel
	case PANIC:
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
