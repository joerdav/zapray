package zapray

import (
	"context"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
)

type ZaprayLogger struct {
	*zap.Logger
}

func NewZaprayLogger(zapLogger *zap.Logger) *ZaprayLogger {
	return &ZaprayLogger{
		Logger: zapLogger,
	}
}

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
