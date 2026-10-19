package repo

import (
	"context"
	"log/slog"
	"singel-sondag/domain"
	"os"

	"github.com/lmittmann/tint"
)

type slogger struct {
	*slog.Logger
}

func NewLogger(env string) domain.Logger {
	return &slogger{
		slog.New(preferredSlogHandler(env)),
	}
}

func preferredSlogHandler(env string) slog.Handler {
	if env == "production" {
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	return tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	})
}

func (s *slogger) Info(ctx context.Context, msg string, args ...any) {
	s.InfoContext(ctx, msg, args...)
}

func (s *slogger) Error(ctx context.Context, msg string, args ...any) {
	s.ErrorContext(ctx, msg, args...)
}
