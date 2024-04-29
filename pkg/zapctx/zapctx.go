package zapctx

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L returns the global logger with considering the Go context.
func L(ctx context.Context) *zap.Logger {
	logger := zap.L()

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return logger
	}

	return logger.With(
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("span_id", span.SpanContext().SpanID().String()),
	)
}

// StartZapctx configure in zap.Globals logs the zapctx logger.
func StartZapctx() error {
	return StartZapctxWithLevel(zapcore.DebugLevel)
}

func StartZapctxWithLevel(logLevel zapcore.Level) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level.SetLevel(logLevel)

	log, err := config.Build()
	if err != nil {
		return err
	}

	_ = zap.ReplaceGlobals(log)
	return nil
}
