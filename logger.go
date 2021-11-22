package zapray

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
)

// Logger is a wrapper for zap.Logger, exposes the zap.Logger functions and adds the ability to Trace logs.
type Logger struct {
	*zap.Logger
}

// NewLogger creates a new instance of *Logger and wraps the logger provided.
//   z, _ := zap.NewProduction()
//   log := zapray.NewLogger(z)
func NewLogger(zapLogger *zap.Logger) *Logger {
	return &Logger{
		Logger: zapLogger,
	}
}

// NewNop creates a new instance of *Logger and includes a zap.NewNop().
//   log := zapray.NewNop()
func NewNop() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
	}
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
func (zprl *Logger) Trace(ctx context.Context) (logger *zap.Logger) {
	logger = zprl.Logger
	defer func() {
		if r := recover(); r != nil {
			zprl.Logger.Warn("no segment found")
		}
	}()
	_, seg := xray.BeginSubsegment(ctx, "zapraylog")
	seg.Close(nil)
	logger = zprl.Logger.With(zap.String("@xrayTraceId", seg.TraceID), zap.String("@xraySegmentId", seg.ParentID))
	seg.ParentSegment.RemoveSubsegment(seg)
	return
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
func (zprl *Logger) TraceRequest(r *http.Request) *zap.Logger {
	return zprl.Trace(r.Context())
}
