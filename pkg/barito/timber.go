package barito

import (
	"os"
	"structs"
	"time"
)

type (
	Timber struct {
		_         structs.HostLayout `json:"-"`
		Message   any                `json:"@message"`
		Metadata  *TimberMetadata    `json:"@metadata,omitempty"`
		Tag       string             `json:"@tag"`
		Timestamp time.Time          `json:"@timestamp"`
	}
	TimberMetadata struct {
		_           structs.HostLayout `json:"-"`
		ExtraLabels map[string]string  `json:"extra_labels,omitempty"`
		Fluentbit   map[string]any     `json:"fluentbit,omitempty"`
		Kubernetes  map[string]any     `json:"kubernetes,omitempty"`
		Logfile     map[string]any     `json:"logfile,omitempty"`
	}
)

func CreateTimber(
	extraLabels map[string]string,
	kubernetesMetadata map[string]any,
	logfileMetadata map[string]any,
	msg any,
	tag string,
	timestamp time.Time,
) *Timber {
	return &Timber{
		Message: msg,
		Metadata: &TimberMetadata{
			ExtraLabels: extraLabels,
			Fluentbit: map[string]any{
				"hostname":          os.Getenv("HOSTNAME"),
				"timber_created_at": time.Now().Format(time.RFC3339Nano),
			},
			Kubernetes: kubernetesMetadata,
			Logfile:    logfileMetadata,
		},
		Tag:       tag,
		Timestamp: timestamp,
	}
}
