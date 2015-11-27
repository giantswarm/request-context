package requestcontext

import (
	"fmt"
	"sync"

	"github.com/op/go-logging"
)

// LoggerRegistry provides the api for a registry of loggers.
type LoggerRegistry interface {
	MustCreate(name string, level ...string) Logger
	Get(name string) (Logger, error)
	List() []string
	GetLevel(name string) (string, error)
	SetLevel(name, level string) error
}

type loggerRegistry struct {
	defaultConfig LoggerConfig
	loggers       map[string]Logger
	mutex         sync.Mutex
}

// NewLoggerRegistry creates and initializes a new LoggerRegistry.
func NewLoggerRegistry(defaultConfig LoggerConfig) LoggerRegistry {
	return &loggerRegistry{
		defaultConfig: defaultConfig,
		loggers:       make(map[string]Logger),
	}
}

// MustCreate creates a new logger and registers it.
func (lg *loggerRegistry) MustCreate(name string, level ...string) Logger {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	if _, ok := lg.loggers[name]; ok {
		panic(fmt.Sprintf("A logger named '%s' already exists", name))
	}

	config := lg.defaultConfig
	config.Name = lg.defaultConfig.Name + "." + name
	if len(level) == 1 {
		config.Level = level[0]
	} else if len(level) > 1 {
		panic("Can only have 1 log level")
	}
	l := MustGetLogger(config)
	lg.loggers[name] = l

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
