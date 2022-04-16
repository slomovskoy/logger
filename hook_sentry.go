package logger

import (
	"log"
	"time"

	"github.com/evalphobia/logrus_sentry"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultHookLevel is a default level, when hook launches
	DefaultHookLevel = "error"
	// SentryHookName is a hook name in the hook list in the config file
	SentryHookName = "sentry"
	// DefaultStackTraceEnable enables/disables stackstrace
	DefaultStackTraceEnable = false
	// DefaultStacktraceLevel is default level to send stacktrace
	DefaultStacktraceLevel = logrus.ErrorLevel
	// DefaultStacktraceContext is a number of lines to send to the setntry
	DefaultStacktraceContext = 3
	// DefaultStacktraceIncludeErrorBreadcrumb whether to create a breadcrumb with the full text of error
	DefaultStacktraceIncludeErrorBreadcrumb = true
	// DefaultStacktraceSendExceptionType - see docs
	DefaultStacktraceSendExceptionType = true
	// DefaultStacktraceSkip count of stack frames to be skipped
	DefaultStacktraceSkip = 6
	// DefaultStacktraceSwitchExceptionTypeAndMessage - see docs
	DefaultStacktraceSwitchExceptionTypeAndMessage = true
)

// SentryHookConfig contains config for setry hook
type SentryHookConfig struct {
	DSN                string              `mapstructure:"dsn"`
	Level              string              `mapstructure:"level"`
	Tags               map[string]string   `mapstructure:"tags"`
	Timeout            int                 `mapstructure:"timeout"`
	StackTraceSettings *StackTraceSettings `mapstructure:"stacktrace_settings"`
}

// StackTraceSettings contains settings for stacktrace
type StackTraceSettings struct {
	Enable                        *bool `mapstructure:"enable"`
	Context                       *int  `mapstructure:"context_lines"`
	IncludeErrorBreadcrumb        *bool `mapstructure:"include_error_breadcrumb"`
	SendExceptionType             *bool `mapstructure:"send_exception_type"`
	Skip                          *int  `mapstructure:"skip_stack_frames"`
	SwitchExceptionTypeAndMessage *bool `mapstructure:"switch_exception_type_and_message"`
}

// GetSetryHookWithSettings prepares nad returns lgorus hook
func GetSetryHookWithSettings(settings interface{}) (logrus.Hook, bool) {
	var setrySettings *SentryHookConfig
	err := mapstructure.Decode(settings, &setrySettings)
	if err != nil {
		log.Printf("init logrus setry hook failed: '%s'", err)
		return nil, false
	}
	if setrySettings.DSN == "" {
		return nil, false
	}
	if setrySettings.Level == "" {
		setrySettings.Level = DefaultHookLevel
	}
	if setrySettings.Tags == nil {
		setrySettings.Tags = make(map[string]string)
	}

	hook, err := logrus_sentry.NewWithTagsSentryHook(
		setrySettings.DSN,
		setrySettings.Tags,
		getlogrusLevels(setrySettings.Level),
	)
	if err != nil {
		log.Printf("init logrus setry hook failed: '%s'", err)
		return nil, false
	}

	hook.Timeout = time.Duration(setrySettings.Timeout) * time.Second
	hook.StacktraceConfiguration.Level = DefaultStacktraceLevel
	hook.StacktraceConfiguration.Enable = DefaultStackTraceEnable
	hook.StacktraceConfiguration.Context = DefaultStacktraceContext
	hook.StacktraceConfiguration.IncludeErrorBreadcrumb = DefaultStacktraceIncludeErrorBreadcrumb
	hook.StacktraceConfiguration.SendExceptionType = DefaultStacktraceSendExceptionType
	hook.StacktraceConfiguration.Skip = DefaultStacktraceSkip
	hook.StacktraceConfiguration.SwitchExceptionTypeAndMessage = DefaultStacktraceSwitchExceptionTypeAndMessage

	if setrySettings.StackTraceSettings != nil {
		hook.StacktraceConfiguration.Level = logrusLevelByString(setrySettings.Level, DefaultStacktraceLevel)
		if setrySettings.StackTraceSettings.Enable != nil {
			hook.StacktraceConfiguration.Enable = *setrySettings.StackTraceSettings.Enable
		}
		if setrySettings.StackTraceSettings.Context != nil {
			hook.StacktraceConfiguration.Context = *setrySettings.StackTraceSettings.Context
		}
		if setrySettings.StackTraceSettings.IncludeErrorBreadcrumb != nil {
			hook.StacktraceConfiguration.IncludeErrorBreadcrumb =
				*setrySettings.StackTraceSettings.IncludeErrorBreadcrumb
		}
		if setrySettings.StackTraceSettings.SendExceptionType != nil {
			hook.StacktraceConfiguration.SendExceptionType = *setrySettings.StackTraceSettings.SendExceptionType
		}
		if setrySettings.StackTraceSettings.Skip != nil {
			hook.StacktraceConfiguration.Skip = *setrySettings.StackTraceSettings.Skip
		}
		if setrySettings.StackTraceSettings.SwitchExceptionTypeAndMessage != nil {
			hook.StacktraceConfiguration.SwitchExceptionTypeAndMessage =
				*setrySettings.StackTraceSettings.SwitchExceptionTypeAndMessage
		}
	}

	return hook, true
}

func getlogrusLevels(lvl string) []logrus.Level {
	var logrusLevel = logrusLevelByString(lvl, DefaultStacktraceLevel)
	var logrusLevels []logrus.Level

	for i := logrusLevel; ; i-- {
		logrusLevels = append(logrusLevels, i)
		if i == logrus.PanicLevel {
			break
		}
	}
	return logrusLevels
}
