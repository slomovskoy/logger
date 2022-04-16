package logger

import "context"

type logCtx struct{}

// WithLogger returns new context with Logger
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logCtx{}, logger)
}

// FromContext returns Logger from context or creates default one
func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(logCtx{}).(Logger)
	if !ok {
		return newWithSettings(
			DefaultInstance, DefaultLogLevel,
			nil, DefaultFormat,
			nil, []Hook{},
		)
	}
	return logger
}
