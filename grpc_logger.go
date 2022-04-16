package logger

import (
	"context"
	"path"
	"time"

	"github.com/google/uuid"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// SystemField is used in every log statement made through logger.
	// Can't be overwritten before any initialization code.
	SystemField = "system"

	// KindField describes the log field used to indicate whether this is a server or a client log statement.
	KindField = "span.kind"

	// TraceID is field for tracing request flow messages
	TraceID = "traceID"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		startTime := time.Now()
		newCtx := newLoggerForCall(ctx, logger, info.FullMethod, startTime)

		resp, err := handler(newCtx, req)

		logFinishedRequest(newCtx, err, startTime)
		return resp, err
	}
}

// StreamServerInterceptor returns grpc stream interceptor with logged context and trace ID
func StreamServerInterceptor(logger Logger) grpc.StreamServerInterceptor {
	return func(
		srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
	) error {
		startTime := time.Now()
		newCtx := newLoggerForCall(stream.Context(), logger, info.FullMethod, startTime)
		wrapped := grpcMiddleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx

		err := handler(srv, wrapped)

		logFinishedRequest(newCtx, err, startTime)
		return err
	}
}

func newLoggerForCall(ctx context.Context, entry Logger, fullMethodString string, start time.Time) context.Context {
	service := path.Dir(fullMethodString)[1:]
	method := path.Base(fullMethodString)
	callLog := entry.WithFields(
		Fields{
			TraceID:           uuid.New().String(),
			SystemField:       "grpc",
			KindField:         "server",
			"grpc.service":    service,
			"grpc.method":     method,
			"grpc.start_time": start.Format(time.RFC3339),
		})

	if d, ok := ctx.Deadline(); ok {
		callLog = callLog.WithFields(
			Fields{
				"grpc.request.deadline": d.Format(time.RFC3339),
			})
	}

	loggedCtx := WithLogger(ctx, callLog)
	return loggedCtx
}

func withCodeLevel(logger Logger, code codes.Code, format string, args ...interface{}) {
	switch code {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.Unauthenticated:
		logger.Infof(format, args...)
	case codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted, codes.OutOfRange,
		codes.Unavailable:
		logger.Warnf(format, args...)
	case codes.Unknown,
		codes.Unimplemented,
		codes.DataLoss,
		codes.Internal:
		logger.Errorf(format, args...)
	default:
		logger.Errorf(format, args...)
	}
}

func logFinishedRequest(ctx context.Context, err error, startTime time.Time) {
	logger := FromContext(ctx)
	st, ok := status.FromError(err)
	fields := Fields{
		"grpc.code":     st.Code(),
		"grpc.message":  st.Message(),
		"grpc.duration": time.Since(startTime),
	}
	if ok {
		withCodeLevel(logger.WithFields(fields), st.Code(), "request finished")
	} else {
		withCodeLevel(logger.WithFields(fields), st.Code(), "request failed")
	}
}
