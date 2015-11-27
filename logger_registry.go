package requestcontext

import (
	"fmt"
	"sync"

	"github.com/op/go-logging"
)

// LoggerRegistry provides the api for a registry of loggers.
type LoggerRegistry interface {
	MustCreate(config LoggerConfig) Logger
	Get(name string) (Logger, error)
	List() []string
	GetLevel(name string) (string, error)
	SetLevel(name, level string) error
}

type loggerRegistry struct {
	loggers map[string]Logger
	mutex   sync.Mutex
}

// NewLoggerRegistry creates and initializes a new LoggerRegistry.
func NewLoggerRegistry() LoggerRegistry {
	return &loggerRegistry{
		loggers: make(map[string]Logger),
	}
}

// MustCreate creates a new logger and registers it.
func (lg *loggerRegistry) MustCreate(config LoggerConfig) Logger {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	if _, ok := lg.loggers[config.Name]; ok {
		panic(fmt.Sprintf("A logger named '%s' already exists", config.Name))
	}

	l := MustGetLogger(config)
	lg.loggers[config.Name] = l
	return l
}

// Get returns the logger with given name or NotFoundError if no such logger
// exists.
func (lg *loggerRegistry) Get(name string) (Logger, error) {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	l, ok := lg.loggers[name]
	if ok {
		return l, nil
	}

	return Logger{}, maskAny(NotFoundError)
}

// List returns a list of names of all registered loggers
func (lg *loggerRegistry) List() []string {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	names := []string{}
	for k, _ := range lg.loggers {
		names = append(names, k)
	}
	return names
}

// GetLevel returns the current log level of the logger with given name
func (lg *loggerRegistry) GetLevel(name string) (string, error) {
	l, err := lg.Get(name)
	if err != nil {
		return "", maskAny(err)
	}
	return l.config.Level, nil
}

// SetLevel changes the log level of the logger with given name
func (lg *loggerRegistry) SetLevel(name, level string) error {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	l, ok := lg.loggers[name]
	if !ok {
		return maskAny(NotFoundError)
	}

	if _, err := logging.LogLevel(level); err != nil {
		return maskAny(err)
	}

	l.config.Level = level
	l.setupBacked(l.config)
	lg.loggers[name] = l

	return nil
}
