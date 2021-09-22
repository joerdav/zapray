package zapray

import (
	"context"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
)

// ZaprayLogger is a wrapper for zap.Logger, exposes the zap.Logger functions and adds the ability to Trace logs.
type ZaprayLogger struct {
	*zap.Logger
}

// NewZaprayLogger creates a new instance of *ZaprayLogger and wraps the logger provided.
//   z, _ := zap.NewProduction()
//   log := zapray.NewZaprayLogger(z)
func NewZaprayLogger(zapLogger *zap.Logger) *ZaprayLogger {
	return &ZaprayLogger{
		Logger: zapLogger,
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
func (zprl *ZaprayLogger) Trace(ctx context.Context) *zap.Logger {
	if zprl == nil {
		return nil
	}
	seg := xray.GetSegment(ctx)
	if seg == nil {
		return zprl.Logger
	}
	traceId := xray.TraceID(ctx)
	return zprl.Logger.With(zap.String("@xrayTraceId", traceId), zap.String("@xraySegmentId", seg.ID))
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
func (zprl *ZaprayLogger) TraceRequest(r *http.Request) *zap.Logger {
	if zprl == nil {
		return nil
	}
	seg := xray.GetSegment(r.Context())
	if seg == nil {
		return zprl.Logger
	}
	traceId := xray.TraceID(r.Context())
	return zprl.Logger.With(zap.String("@xrayTraceId", traceId), zap.String("@xraySegmentId", seg.ID))
}
