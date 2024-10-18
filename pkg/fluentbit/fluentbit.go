package fluentbit

import (
	"log/slog"
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
		logs.Err(err)
		return output.FLB_ERROR
	}

	output.FLBPluginSetContext(plugin, configuration)

	logs.Debug(PluginName+" output plugin initialized",
		slog.String("version", PluginVersion),
		slog.String("build_time", PluginBuildTime),
		configuration.ToSlogAttr(),
	)
	return output.FLB_OK
}

func ParseRecordTimestamp(ts any) time.Time {
	switch t := ts.(type) {
	case output.FLBTime:
		return ts.(output.FLBTime).Time
	case uint64:
		return time.Unix(int64(t), 0)
	default:
		logs.Debug("time provided is invalid, creating one",
			slog.Any("provided_timestamp", ts),
		)
		return time.Now()
	}
}
