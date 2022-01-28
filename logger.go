package zapray

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AppVersion is exposed to allow build time versions to be set.
// go build -ldflags="-X 'github.com/joe-davidson1802/zapray.AppVersion.AppVersion=v1.0.0'"
var AppVersion = ""

// Logger is a wrapper for zap.Logger, exposes the zap.Logger functions and adds the ability to Trace logs.
type Logger struct {
	*zap.Logger
	SuppressMissingSegmentWarning bool
}

// NewLogger creates a new instance of *Logger and wraps the logger provided.
//   z, _ := zap.NewProduction()
//   log := zapray.NewLogger(z)
func NewLogger(zapLogger *zap.Logger) *Logger {
	l := zapLogger
	if AppVersion != "" {
		l = zapLogger.With(zap.String("appVersion", AppVersion))
	}
	return &Logger{
		Logger: l,
	}
}

// New creates a new instance of zap.Logger and wraps it in a zapray.Logger
func New(core zapcore.Core, options ...zap.Option) *Logger {
	return NewLogger(zap.New(core, options...))
}

// NewNop creates a new instance of *Logger and includes a zap.NewNop().
//   log := zapray.NewNop()
func NewNop() *Logger {
	return NewLogger(zap.NewNop())
}

// NewDevelopment creates a new instance of *Logger and includes a zap.NewDevelopment().
//   log := zapray.NewDevelopment()
func NewDevelopment() (*Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return NewLogger(logger), nil
}

// NewProduction creates a new instance of *Logger and includes a zap.NewProduction().
//   log := zapray.NewProduction()
func NewProduction() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return NewLogger(logger), nil
}

// Trace creates a new zap.Logger but with the xrayTraceId and xraySegmentId baked in.
//   log.Trace(ctx).Info("myLog")
//
// Once trace is called you can use it as a zap.Logger.
//
//   tracedLogger := log.Trace(ctx)
//   tracedLogger.Info("Log one")
//   tracedLogger.Info("Log two")
//
// This means as above you can trace once and use the provided logger.
func (zprl *Logger) Trace(ctx context.Context) (res *Logger) {
	logger := zprl.Logger
	defer func() {
		if r := recover(); r != nil && !zprl.SuppressMissingSegmentWarning {
			zprl.Logger.Warn("no segment found")
		}
		res = zprl
	}()
	seg := xray.GetSegment(ctx)
	traceId := seg.TraceID
	segmentId := seg.ID
	logger = zprl.Logger.With(zap.String("@xrayTraceId", traceId), zap.String("@xraySegmentId", segmentId))
	return NewLogger(logger)
}

// TraceRequest creates a new zap.Logger but with the xrayTraceId and xraySegmentId baked in.
//   log.TraceRequest(r).Info("myLog")
//
// Once trace is called you can use it as a zap.Logger.
//
//   tracedLogger := log.Trace(r)
//   tracedLogger.Info("Log one")
//   tracedLogger.Info("Log two")
//
// This means as above you can trace once and use the provided logger.
func (zprl *Logger) TraceRequest(r *http.Request) *Logger {
	return zprl.Trace(r.Context())
}

// WithOptions delegates to zap.Logger.WithOptions and wraps the resulting logger into a zapray.Logger
func (log *Logger) WithOptions(opts ...zap.Option) *Logger {
	return NewLogger(log.Logger.WithOptions(opts...))
}

// With delegates to zap.Logger.With and wraps the resulting logger into a zapray.Logger
func (log *Logger) With(fields ...zap.Field) *Logger {
	return NewLogger(log.Logger.With(fields...))
}

// Named delegates to zap.Logger.Named and wraps the resulting logger into a zapray.Logger
func (log *Logger) Named(s string) *Logger {
	return NewLogger(log.Logger.Named(s))
}
