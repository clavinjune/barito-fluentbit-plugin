package barito

import (
	"os"
	"structs"
	"time"
)

type (
	Timber struct {
		_     structs.HostLayout `json:"-"`
		Items []*Item            `json:"items"`
	}

	Item struct {
		_ structs.HostLayout `json:"-"`
		// ClienTrail stores client data such as hostname
		ClientTrail map[string]any `json:"@client_trail"`
		// Message comes from the fluentbit Record
		Message map[string]any `json:"@message"`
		// Tag comses from the fluentbit configuration
		Tag string `json:"@tag"`
		// Timestamp comes from the fluentbit record
		Timestamp time.Time `json:"@timestamp"`
	}
)

// createTimber constructs Timber with multiple messages
func createTimber(tag string, timestamp time.Time, msgs ...map[string]any) *Timber {
	items := make([]*Item, 0, len(msgs))
	for _, msg := range msgs {
		items = append(items, &Item{
			ClientTrail: map[string]any{
				"hostname": os.Getenv("HOSTNAME"),
			},
			Message:   msg,
			Tag:       tag,
			Timestamp: timestamp,
		})
	}

	return &Timber{Items: items}
}
