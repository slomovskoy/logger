package logger

import "github.com/sirupsen/logrus"

// Config prototype
type Config struct {
	LogLevel          string `mapstructure:"log_level"`
	FormatterName     string `mapstructure:"formatter_name"`
	FormatterInstance logrus.Formatter
	Hooks             Hooks `mapstructure:"hooks"`
}
