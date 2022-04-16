package logger

import "github.com/sirupsen/logrus"

// logWrapper is a wrapper around *logrus.Entry
type logWrapper struct {
	*logrus.Entry
}

func (lw *logWrapper) WithError(err error) Logger {
	return &logWrapper{lw.Entry.WithError(err)}
}

func (lw *logWrapper) WithField(key string, value interface{}) Logger {
	return &logWrapper{lw.Entry.WithField(key, value)}
}

func (lw *logWrapper) WithFields(fields Fields) Logger {
	return &logWrapper{lw.Entry.WithFields(logrus.Fields(fields))}
}
