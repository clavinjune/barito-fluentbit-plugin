package barito

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"structs"
)

type (
	Client struct {
		_      structs.HostLayout
		client *http.Client
		Config *Configuration
	}
	ProduceBatchRequest struct {
		_     structs.HostLayout `json:"-"`
		Items []*Timber          `json:"items"`
	}
)

func NewClient(client *http.Client, config *Configuration) *Client {
	return &Client{
		client: client,
		Config: config,
	}
}

func (c *Client) ProduceBatch(ctx context.Context, r *ProduceBatchRequest) error {
	url := fmt.Sprintf("%s/produce_batch", c.Config.BaritoHost)
	ctx, cancel := context.WithTimeout(ctx, c.Config.BaritoRequestTimeoutDuration)
	defer cancel()

	req, err := c.createCompressedRequest(ctx, url, r)
	if err != nil {
		return fmt.Errorf("barito: failed when creating http request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("barito: failed when sending http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("barito: response error: %v, error when reading response: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("barito: response error: %v, body: %v", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) createCompressedRequest(ctx context.Context, url string, r *ProduceBatchRequest) (*http.Request, error) {
	pr, pw := io.Pipe()
	go func(ctx context.Context, request *ProduceBatchRequest) {
		gw := gzip.NewWriter(pw)
		if err := json.NewEncoder(gw).Encode(request); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
		if err := gw.Close(); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
		if err := pw.Close(); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
	}(ctx, r)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, pr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Name", c.Config.ApplicationName)
	req.Header.Set("X-App-Group-Secret", c.Config.ApplicationGroupSecret)
	return req, nil
}
