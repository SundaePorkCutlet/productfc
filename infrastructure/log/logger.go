package log

import (
	// golang package
	"context"
	"os"
	"time"

	// external package
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var Logger *zerolog.Logger

// SetupLogger setup logger.
func SetupLogger() {
	// logrus와 동일한 스타일의 텍스트 포맷터 설정
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false, // ForceColors: true와 동일
	}

	logger := zerolog.New(output).With().
		Timestamp().
		Logger()

	Logger = &logger

	// logrus와 동일한 초기화 메시지
	Logger.Info().Msg("Logged initiated using zerolog!")
}

// LogWithTrace log with trace.
//
// It returns pointer of zerolog.Event when successful.
// Otherwise, nil pointer of zerolog.Event will be returned.
func LogWithTrace(ctx context.Context) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID().String()
	return Logger.Info().Str("trace_id", traceID)
}
