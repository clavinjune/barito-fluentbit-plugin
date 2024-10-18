package logs

import (
	"log/slog"
	"strconv"
	"structs"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

type (
	Configuration struct {
		_       structs.HostLayout `json:"-"`
		IsDebug bool               `json:"is_debug"`
		IsJSON  bool               `json:"is_json"`
	}
)

func (c *Configuration) ToSlogAttr() slog.Attr {
	return slog.Group("log",
		slog.Bool("is_debug", c.IsDebug),
		slog.Bool("is_json", c.IsJSON),
	)
}

func GetConfigurationFromPlugin(plugin unsafe.Pointer) *Configuration {
	isDebug, err := strconv.ParseBool(output.FLBPluginConfigKey(plugin, "is_debug"))
	if err != nil {
		isDebug = false
	}

	isJSON, err := strconv.ParseBool(output.FLBPluginConfigKey(plugin, "is_json"))
	if err != nil {
		isJSON = false
	}

	return &Configuration{
		IsDebug: isDebug,
		IsJSON:  isJSON,
	}
}
