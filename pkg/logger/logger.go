package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddlware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (it *HandlerMiddlware) Enabled(ctx context.Context, level slog.Level) bool {
	if os.Getenv("APP_ENV") == "DEV" {
		return level >= slog.LevelDebug
	}
	return level >= slog.LevelInfo
}

func (it *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	// Будем добавлять неизменяющиеся данные юзер роль и т.п
	return it.next.Handle(ctx, rec)
}

func (it *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: it.next.WithAttrs(attrs)}
}

func (it *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: it.next.WithGroup(name)}
}

func InitLogging() *slog.Logger {
	appEnv := os.Getenv("APP_ENV")

	opts := &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("2006-01-02T15:04:05"))
			}

			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				// Создаем вложенный объект для source
				return slog.Attr{
					Key: "called_at",
					Value: slog.GroupValue(
						slog.String("method", filepath.Base(source.Function)),
						slog.String("file", source.File),
						slog.Int("line", source.Line),
					),
				}
			}

			return a
		},
	}

	var handler slog.Handler
	if appEnv == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	handler = NewHandlerMiddlware(handler)

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
