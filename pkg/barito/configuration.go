package barito

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"structs"
	"time"
	"unsafe"

	"github.com/clavinjune/barito-fluentbit-plugin/pkg/logs"
	"github.com/fluent/fluent-bit-go/output"
)

var (
	defaultBaritoRequestTimeoutDuration time.Duration = time.Minute
	errInvalidURL                       error         = errors.New("url is missing scheme or port")
)

type (
	Configuration struct {
		_                            structs.HostLayout  `json:"-"`
		ApplicationName              string              `json:"application_name"`
		ApplicationGroupSecret       string              `json:"application_group_secret"`
		BaritoHost                   string              `json:"barito_host"`
		BaritoRequestTimeoutDuration time.Duration       `json:"barito_request_timeout_duration"`
		ExtraLabels                  string              `json:"extra_labels"`
		LogConfiguration             *logs.Configuration `json:"log_configuration"`
		ParsedExtraLabels            map[string]string   `json:"parsed_extra_labels"`
	}
)

func (c *Configuration) ToSlogAttr() slog.Attr {
	return slog.Group("barito",
		slog.String("application_name", c.ApplicationName),
		slog.String("application_group_secret", c.ApplicationGroupSecret),
		slog.String("barito_host", c.BaritoHost),
		slog.Duration("barito_request_timeout_duration", c.BaritoRequestTimeoutDuration),
		slog.Any("extra_labels", c.ParsedExtraLabels),
		c.LogConfiguration.ToSlogAttr(),
	)
}

func GetConfigurationFromPlugin(plugin unsafe.Pointer) (*Configuration, error) {
	appName := strings.TrimSpace(output.FLBPluginConfigKey(plugin, "application_name"))
	if appName == "" {
		return nil, errors.New("barito: application_name must be filled")
	}

	appGroupSecret := strings.TrimSpace(output.FLBPluginConfigKey(plugin, "application_group_secret"))
	if appGroupSecret == "" {
		return nil, errors.New("barito: application_group_secret must be filled")
	}

	baritoHost, err := parseURL(output.FLBPluginConfigKey(plugin, "barito_host"))
	if err != nil {
		return nil, fmt.Errorf("barito: barito_host is invalid: %w", err)
	}

	baritoRequestTimeoutDurationStr := strings.TrimSpace(output.FLBPluginConfigKey(plugin, "barito_request_timeout_duration"))
	baritoRequestTimeoutDuration, err := time.ParseDuration(baritoRequestTimeoutDurationStr)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelWarn, "barito: barito_request_timeout_duration is invalid, using default timeout",
			slog.String("invalid_duration", baritoRequestTimeoutDurationStr),
			slog.Duration("default_timeout", defaultBaritoRequestTimeoutDuration),
		)
		baritoRequestTimeoutDuration = defaultBaritoRequestTimeoutDuration
	}

	extraLabels := strings.TrimSpace(output.FLBPluginConfigKey(plugin, "extra_labels"))

	return &Configuration{
		ApplicationName:              appName,
		ApplicationGroupSecret:       appGroupSecret,
		BaritoHost:                   baritoHost,
		BaritoRequestTimeoutDuration: baritoRequestTimeoutDuration,
		ExtraLabels:                  extraLabels,
		LogConfiguration:             logs.GetConfigurationFromPlugin(plugin),
		ParsedExtraLabels:            parseExtraLabels(extraLabels),
	}, nil
}

func parseURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errInvalidURL
	}
	return parsedURL.String(), nil
}

func parseExtraLabels(el string) map[string]string {
	m := make(map[string]string)
	for _, label := range strings.Split(el, ",") {
		splittedLabel := strings.Split(label, "=")
		if len(splittedLabel) != 2 {
			continue
		}
		m[strings.TrimSpace(splittedLabel[0])] = strings.TrimSpace(splittedLabel[1])
	}

	return m
}
