package fluentbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"unsafe"

	"github.com/clavinjune/barito-fluentbit-plugin/pkg/barito"
	"github.com/clavinjune/barito-fluentbit-plugin/pkg/logs"
	"github.com/fluent/fluent-bit-go/output"
)

func PluginInit(
	PluginName,
	PluginVersion,
	PluginBuildTime string,
	plugin unsafe.Pointer,
) int {
	logs.SetDefaultLogger(logs.GetConfigurationFromPlugin(plugin))

	configuration, err := barito.GetConfigurationFromPlugin(plugin)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, err.Error())
		return output.FLB_ERROR
	}

	baritoClient := barito.NewClient(&http.Client{}, configuration)

	output.FLBPluginSetContext(plugin, baritoClient)

	slog.LogAttrs(context.Background(), slog.LevelDebug, PluginName+" output plugin initialized",
		slog.String("version", PluginVersion),
		slog.String("build_time", PluginBuildTime),
		configuration.ToSlogAttr(),
	)
	return output.FLB_OK
}

func ParseRecordData(rec any) any {
	switch t := rec.(type) {
	case []uint8:
		m := make(map[string]any)
		if err := json.Unmarshal(t, &m); err != nil {
			return string(t)
		}
		return m
	case nil:
		return ""
	default:
		slog.LogAttrs(context.Background(), slog.LevelDebug, "failed parsing record data",
			slog.String("type", fmt.Sprintf("%T", t)),
			slog.String("content", fmt.Sprintf("%v", t)),
		)
		return fmt.Sprintf("%v", t)
	}
}

func ParseRecordTimestamp(ts any) time.Time {
	switch t := ts.(type) {
	case output.FLBTime:
		return ts.(output.FLBTime).Time
	case uint64:
		return time.Unix(int64(t), 0)
	default:
		slog.LogAttrs(context.Background(), slog.LevelDebug, "time provided is invalid, creating one",
			slog.Any("provided_timestamp", ts),
		)
		return time.Now()
	}
}
