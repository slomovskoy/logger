package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	// StartSkip is a start offset for runtime.Caller()
	StartSkip int = 2
	// FileLineHookName is a filename line hook name
	FileLineHookName string = "filename_line"
	// DefaultFileNameLineKey is a default field name in log string
	DefaultFileNameLineKey = "where"
)

// FileLineHook contains caller's log settings
type FileLineHook struct {
	LogKeyName string `mapstructure:"field_name"`
}

// GetFileLineHookWithSettings prepares and returns filename line hook
func GetFileLineHookWithSettings(settings interface{}) (*FileLineHook, bool) {
	var hook *FileLineHook
	err := mapstructure.Decode(settings, &hook)
	if err != nil {
		logrus.Printf("init file line hook failed: '%s'", err)
		return nil, false
	}
	if hook.LogKeyName == "" {
		hook.LogKeyName = DefaultFileNameLineKey
	}
	return hook, true
}

// Levels implements logrus's Hook interface
func (hook *FileLineHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implements logrus's Hook interface
func (hook *FileLineHook) Fire(entry *logrus.Entry) error {
	var (
		file string
		line int
	)
	for i := 0; i < 10; i++ {
		file, line = getCaller(StartSkip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}

	entry.Data[hook.LogKeyName] = fmt.Sprintf("%s:%d", file, line)
	return nil
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}

	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}

	return file, line
}
