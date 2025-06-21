package context

import (
	"context"
	"log/slog"
)

type key struct{}

// WithLogger добавляет логгер в контекст
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

// LoggerFromContext возвращает логгер из контекста
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(key{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
