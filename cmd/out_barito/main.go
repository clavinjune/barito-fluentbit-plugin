package main

import (
	"C"
	"context"
	"log/slog"
	"time"
	"unsafe"

	"github.com/clavinjune/barito-fluentbit-plugin/pkg/barito"
	"github.com/clavinjune/barito-fluentbit-plugin/pkg/fluentbit"
	"github.com/fluent/fluent-bit-go/output"
)

var (
	PluginVersion   = "dev"
	PluginBuildTime = "N/A"
)

const (
	PluginName        = "barito"
	PluginDescription = "Barito output plugin for Fluent Bit"
)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	return output.FLBPluginRegister(
		def,
		PluginName,
		PluginDescription,
	)
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	return fluentbit.PluginInit(
		PluginName,
		PluginVersion,
		PluginBuildTime,
		plugin,
	)
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {
	n := time.Now()

	slog.LogAttrs(
		context.Background(),
		slog.LevelDebug,
		"start flushing",
	)
	defer func(t time.Time) {
		slog.LogAttrs(
			context.Background(),
			slog.LevelDebug,
			"end flushing",
			slog.Duration("duration", time.Since(t)),
		)
	}(n)

	baritoClient := output.FLBPluginGetContext(ctx).(*barito.Client)

	d := output.NewDecoder(data, int(length))
	timbers := make([]*barito.Timber, 0)
	for {
		ret, ts, rec := output.GetRecord(d)
		if ret != 0 {
			break
		}

		parsedTs := fluentbit.ParseRecordTimestamp(ts)

		logfileMetadata := make(map[string]any)
		kubernetesMetadata := make(map[string]any)
		var msg any

		for k, v := range rec {
			data := fluentbit.ParseRecordData(v)
			kstr := k.(string)

			switch kstr {
			case "log":
				msg = data
			case "kubernetes":
				if d, ok := data.(map[string]any); ok {
					kubernetesMetadata = d
				} else {
					kubernetesMetadata = nil
				}
			default:
				logfileMetadata[kstr] = data
			}
		}
		timbers = append(timbers,
			barito.CreateTimber(
				baritoClient.Config.ParsedExtraLabels,
				kubernetesMetadata,
				logfileMetadata,
				msg,
				C.GoString(tag),
				parsedTs,
			))
	}
	if err := baritoClient.ProduceBatch(context.Background(), &barito.ProduceBatchRequest{
		Items: timbers,
	}); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, err.Error(), baritoClient.Config.ToSlogAttr())
		return output.FLB_RETRY
	}

	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

//export FLBPluginUnregister
func FLBPluginUnregister(def unsafe.Pointer) {
	output.FLBPluginUnregister(def)
}

func main() {
}
