package barito

import (
	"os"
	"structs"
	"time"
)

type (
	Timber struct {
		_ structs.HostLayout `json:"-"`
		// ClienTrail stores client data such as hostname
		ClientTrail map[string]any    `json:"@client_trail"`
		ExtraLabels map[string]string `json:"@extra_labels"`
		// Message comes from the fluentbit Record
		Message map[string]any `json:"@message"`
		// Tag comes from the fluentbit configuration
		Tag string `json:"@tag"`
		// Timestamp comes from the fluentbit record
		Timestamp time.Time `json:"@timestamp"`
	}
)

func CreateTimber(tag string, timestamp time.Time, extraLabels map[string]string, msg map[string]any) *Timber {
	return &Timber{
		ClientTrail: map[string]any{
			"hostname":          os.Getenv("HOSTNAME"),
			"timber_created_at": time.Now().Format(time.RFC3339Nano),
		},
		ExtraLabels: extraLabels,
		Message:     msg,
		Tag:         tag,
		Timestamp:   timestamp,
	}
}
