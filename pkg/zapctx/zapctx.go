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
func StartZapctx(version string) error {
	return StartZapctxWithLevel(version, zapcore.DebugLevel)
}

func StartZapctxWithLevel(version string, logLevel zapcore.Level) error {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level.SetLevel(logLevel)

	log, err := config.Build(zap.Fields(zap.String("version", version)))
	if err != nil {
		return err
	}

	_ = zap.ReplaceGlobals(log)
	return nil
}
