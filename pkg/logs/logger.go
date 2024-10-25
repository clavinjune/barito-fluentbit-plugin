package logs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var (
	redactedLogValue slog.Value = slog.StringValue("[REDACTED]")
)

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
