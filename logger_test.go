package zapray

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-xray-sdk-go/xray"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZaprayLogger(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	log := zap.New(core)
	ctx := context.Background()
	segCtx, seg := xray.BeginSegment(ctx, "testseg")
	defer seg.Close(nil)

	zprl := NewLogger(log)
	zprl.Trace(segCtx).Info("test log")

	if len(recorded.All()) != 1 {
		t.Errorf("expected 1 call got %d", len(recorded.All()))
	}
	containsTrace := ""
	for _, f := range recorded.All()[0].Context {
		if f.Key == "@xrayTraceId" {
			containsTrace = f.String
		}
	}
	fmt.Println("seg", seg.ID)
	fmt.Println("seg", xray.TraceID(segCtx))
	if containsTrace == "" {
		t.Errorf("expected log to contain traceId")
	}

}
