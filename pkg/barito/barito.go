package barito

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	client *http.Client = &http.Client{}
)

func Flush(ctx context.Context, c *Configuration, tag string, timestamp time.Time, msgs ...map[string]any) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	req, err := createCompressedRequest(ctx, c, createTimber(tag, timestamp, msgs...))
	if err != nil {
		return fmt.Errorf("barito: failed when creating http request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("barito: failed when sending http request: %w", err)
	}
	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("barito: error when reading response: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelDebug, string(respB))
	return nil
}

func createCompressedRequest(ctx context.Context, c *Configuration, t *Timber) (*http.Request, error) {
	pr, pw := io.Pipe()
	go func(ctx context.Context, timber *Timber) {
		gw := gzip.NewWriter(pw)
		if err := json.NewEncoder(gw).Encode(timber); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
		if err := gw.Close(); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
		if err := pw.Close(); err != nil {
			slog.LogAttrs(ctx, slog.LevelError, err.Error())
		}
	}(ctx, t)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.ProduceURL, pr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Name", c.ApplicationName)
	req.Header.Set("X-App-Group-Secret", c.ApplicationGroupSecret)
	return req, nil
}
