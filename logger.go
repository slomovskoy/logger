package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// EnvName is an ENV var's name for application environment
	EnvName string = "GO_ENV"
	// EnvLevel is an ENV var's name for logging level
	EnvLevel string = "LOG_LEVEL"
	// DefaultLogLevel is a default log level if it is not passed whether config, nor ENV
	DefaultLogLevel = "info"
	// DefaultInstance is a default instance type
	DefaultInstance = "development"
	// DefaultFormat is a default output messages format
	DefaultFormat = "text"
)

// Fields is a wrapper around logrus.Field
type Fields logrus.Fields

// Logger interface
type Logger interface {
	Trace(...interface{})
	Tracef(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Warning(...interface{})
	Warningf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
}

// New creates a new logger instance (constructor)
func New(conf Config, baseFields map[string]string) Logger {
	loggerFields := make(map[string]interface{})
	for k, v := range baseFields {
		loggerFields[k] = v
	}

	logLevel := conf.LogLevel // get log level from config
	envLogLevel := os.Getenv(EnvLevel)

	if envLogLevel != "" { // if env defined
		logLevel = envLogLevel // get log level from env
	}

	return newWithSettings(
		os.Getenv(EnvName), logLevel,
		conf.FormatterInstance, conf.FormatterName,
		loggerFields, conf.Hooks,
	)
}

func newWithSettings(instance, logLevel string,
	formatter logrus.Formatter,
	formatterName string,
	baseFields map[string]interface{},
	hooks Hooks) Logger {
	switch {
	case instance == "production":
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: false,
			TimestampFormat:  time.RFC1123,
		})
	case formatter != nil:
		logrus.SetFormatter(formatter)
	case formatterName != "":
		if strings.EqualFold(formatterName, "custom") {
			logrus.SetFormatter(new(CustomFormatter))
		}
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC1123,
		})
	}

	var configHooks []Hook
	switch typedHooks := hooks.(type) {
	case []interface{}:
		configHooks = listToHooks(typedHooks)
	case interface{}:
		mapHooks, ok := typedHooks.(map[string]interface{})
		if ok {
			for k, v := range mapHooks {
				configHooks = append(configHooks, Hook{
					Name:     k,
					Settings: v,
				})
			}
		}
	default:
		// no hooks
		configHooks = []Hook{}
	}

	for i := 0; i < len(configHooks); i++ {
		hook := configHooks[i]
		var (
			logrusHook logrus.Hook
			ok         bool
		)
		switch hook.Name {
		case FileLineHookName:
			logrusHook, ok = GetFileLineHookWithSettings(hook.Settings)
		case SentryHookName:
			logrusHook, ok = GetSetryHookWithSettings(hook.Settings)
		}
		if ok { // it is logrus hook
			logrus.AddHook(logrusHook)
		}
	}

	logrus.SetOutput(os.Stdout)
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}

	logrusLevel := logrusLevelByString(strings.ToLower(logLevel), logrus.InfoLevel)
	logrus.SetLevel(logrusLevel)
	return &logWrapper{logrus.WithFields(baseFields)}
}

func logrusLevelByString(level string, defaultLevel logrus.Level) logrus.Level {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrusLevel = defaultLevel
	}
	return logrusLevel
}
