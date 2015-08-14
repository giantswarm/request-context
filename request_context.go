package requestcontext

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/juju/errgo"
	"github.com/op/go-logging"
)

const (
	separator = " | "
)

type Ctx map[string]interface{}

func (ctx Ctx) isEmpty() bool {
	return len(ctx) == 0
}

type Logger struct {
	config LoggerConfig
	logger *logging.Logger
}

type LoggerConfig struct {
	Name  string
	Level string
	Color bool
}

// NewSimpleLogger creates a new logger with a default backend logging to `os.Stderr`.
func MustGetLogger(config LoggerConfig) Logger {
	logger := Logger{
		config: config,
		logger: logging.MustGetLogger(config.Name),
	}

	// See https://godoc.org/github.com/op/go-logging#NewStringFormatter for format verbs.
	format := strings.Join([]string{
		"%{time:2006-01-02 15:04:05}",
		"%{level}",
		"%{message}",
	}, separator)

	formatter := logging.MustStringFormatter(format)
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backend.Color = config.Color
	backendFormatter := logging.NewBackendFormatter(backend, formatter)
	leveledBackend := logging.AddModuleLevel(backendFormatter)

	// Set log level.
	if config.Level != "" {
		logLevel, err := logging.LogLevel(config.Level)
		if err != nil {
			panic(errgo.Mask(err))
		}

		leveledBackend.SetLevel(logLevel, config.Name)
	}

	logger.logger.SetBackend(leveledBackend)

	return logger
}

func (l Logger) Critical(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.CRITICAL) {
		return
	}

	l.logger.Critical(l.extendFormat(ctx, f), v...)
}

func (l Logger) Error(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.ERROR) {
		return
	}

	l.logger.Error(l.extendFormat(ctx, f), v...)
}

func (l Logger) Warning(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.WARNING) {
		return
	}

	l.logger.Warning(l.extendFormat(ctx, f), v...)
}

func (l Logger) Notice(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.NOTICE) {
		return
	}

	l.logger.Notice(l.extendFormat(ctx, f), v...)
}

func (l Logger) Info(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.INFO) {
		return
	}

	l.logger.Info(l.extendFormat(ctx, f), v...)
}

func (l Logger) Debug(ctx Ctx, f string, v ...interface{}) {
	if !l.logger.IsEnabledFor(logging.DEBUG) {
		return
	}

	l.logger.Debug(l.extendFormat(ctx, f), v...)
}

func (l Logger) extendFormat(ctx Ctx, f string) string {
	meta := ""
	if !ctx.isEmpty() {
		rawMeta, err := json.Marshal(ctx)
		if err != nil {
			panic(errgo.Mask(err))
		}
		meta = string(rawMeta)
	}

	format := f
	if meta != "" {
		format += separator + meta
	}

	return format
}
