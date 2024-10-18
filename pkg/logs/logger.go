package logs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	redactedLogValue slog.Value = slog.StringValue("[REDACTED]")
)

func Debug(msg string, attrs ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		msg,
		attrs...,
	)
}

func Info(msg string, attrs ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg,
		attrs...,
	)
}

func Warn(msg string, attrs ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelWarn,
		msg,
		attrs...,
	)
}

func Err(err error, attrs ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelError,
		err.Error(),
		attrs...,
	)
}

func ErrMsg(msg string, attrs ...slog.Attr) {
	slog.LogAttrs(
		context.Background(),
		slog.LevelError,
		msg,
		attrs...,
	)
}

func SetDefaultLogger(c *Configuration) {
	var l slog.Leveler
	if c.IsDebug {
		l = slog.LevelDebug
	} else {
		l = slog.LevelInfo
	}

	ho := &slog.HandlerOptions{
		AddSource: true,
		Level:     l,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				src := a.Value.Any().(*slog.Source)
				return slog.Attr{Key: a.Key, Value: slog.StringValue(fmt.Sprintf("%s:%d", src.File, src.Line))}
			}

			key := strings.ToLower(a.Key)
			if strings.Contains(key, "secret") ||
				strings.Contains(key, "token") {
				return slog.Attr{Key: a.Key, Value: redactedLogValue}
			}

			return a
		},
	}

	var h slog.Handler
	if c.IsJSON {
		h = slog.NewJSONHandler(os.Stdout, ho)
	} else {
		h = slog.NewTextHandler(os.Stdout, ho)
	}

	slog.SetDefault(slog.New(h))
}

func TrackDuration(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) func() {
	n := time.Now()

	slog.LogAttrs(
		ctx,
		level,
		"start "+msg,
		attrs...,
	)

	return func() {
		slog.LogAttrs(
			ctx,
			level,
			"end "+msg,
			append(attrs, slog.Duration("duration", time.Since(n)))...,
		)
	}
}
