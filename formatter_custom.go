package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// CustomFormatter is a container for custom log format
type CustomFormatter struct{}

// Format implements logrus's Formatter interface
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s [%s] <%s>: %s\n", entry.Time, entry.Level, entry.Data, entry.Message)), nil
}
