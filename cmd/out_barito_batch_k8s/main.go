package main

import (
	"C"
	"context"
	"log/slog"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"

	"github.com/clavinjune/barito-fluentbit-plugin/pkg/barito"
	"github.com/clavinjune/barito-fluentbit-plugin/pkg/fluentbit"
)
import "time"

var (
	PluginVersion   = "dev"
	PluginBuildTime = "N/A"
)

const (
	PluginName        = "barito_batch_k8s"
	PluginDescription = "Barito Batch Kubernetes output plugin for Fluent Bit"
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

	configuration := output.FLBPluginGetContext(ctx).(*barito.Configuration)
	d := output.NewDecoder(data, int(length))
	for {
		ret, ts, rec := output.GetRecord(d)
		if ret != 0 {
			break
		}

		parsedTs := fluentbit.ParseRecordTimestamp(ts)

		msgs := make([]map[string]any, 0, len(rec))
		for k, v := range rec {
			msgs = append(msgs, map[string]any{
				k.(string): v,
			})
		}
		if err := barito.Flush(context.Background(), configuration, C.GoString(tag), parsedTs, msgs...); err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, err.Error())
			return output.FLB_RETRY
		}
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
