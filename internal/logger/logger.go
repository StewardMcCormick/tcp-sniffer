package logger

import (
	"io"
	"log/slog"
	"strings"

	"github.com/StewardMcCormick/tcp-analyzer/internal/config"
)

func New(cfg config.LoggerConfig, output io.Writer) *slog.Logger {
	var level slog.LevelVar
	level.Set(parseLevel(cfg.Level))

	opts := &slog.HandlerOptions{
		Level:     &level,
		AddSource: false,
	}

	var handler slog.Handler

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
