package barito

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"structs"
	"unsafe"

	"github.com/clavinjune/barito-fluentbit-plugin/pkg/logs"
	"github.com/fluent/fluent-bit-go/output"
)

type (
	Configuration struct {
		_                      structs.HostLayout  `json:"-"`
		ApplicationName        string              `json:"application_name"`
		ApplicationGroupSecret string              `json:"application_group_secret"`
		ClusterName            string              `json:"cluster_name"`
		ProduceURL             string              `json:"produce_url"`
		LogConfiguration       *logs.Configuration `json:"log_configuration"`
	}
)

func (c *Configuration) ToSlogAttr() slog.Attr {
	return slog.Group("barito",
		slog.String("application_name", c.ApplicationName),
		slog.String("application_group_secret", c.ApplicationGroupSecret),
		slog.String("cluster_name", c.ClusterName),
		slog.String("produce_url", c.ProduceURL),
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

	produceURL, err := parseURL(output.FLBPluginConfigKey(plugin, "produce_url"))
	if err != nil {
		return nil, fmt.Errorf("barito: produce_url is invalid: %w", err)
	}

	return &Configuration{
		ApplicationName:        appName,
		ApplicationGroupSecret: appGroupSecret,
		ClusterName:            strings.TrimSpace(output.FLBPluginConfigKey(plugin, "cluster_name")),
		ProduceURL:             produceURL,
		LogConfiguration:       logs.GetConfigurationFromPlugin(plugin),
	}, nil
}

func parseURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("url is missing scheme or port")
	}

	return parsedURL.String(), nil
}
