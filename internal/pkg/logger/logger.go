package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type LokiWriter struct {
	url    string
	labels map[string]string
	client *http.Client
	mu     sync.Mutex
}

func NewLokiWriter(url string, labels map[string]string) *LokiWriter {
	return &LokiWriter{
		url:    url,
		labels: labels,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (lw *LokiWriter) Write(p []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	entry := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": lw.labels,
				"values": [][]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), string(p)},
				},
			},
		},
	}

	body, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", lw.url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := lw.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return 0, fmt.Errorf("loki responded with status %d", resp.StatusCode)
	}

	return len(p), nil
}

func New(level, filePath, lokiURL string, lokiLabels map[string]string) *slog.Logger {
	opts := &slog.HandlerOptions{}

	switch level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "info":
		opts.Level = slog.LevelInfo
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if filePath != "" {
		if err := os.MkdirAll("logs", 0755); err == nil {
			f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				writers = append(writers, f)
			}
		}
	}

	if lokiURL != "" {
		lw := NewLokiWriter(lokiURL, lokiLabels)
		writers = append(writers, lw)
	}

	multiWriter := io.MultiWriter(writers...)
	handler := slog.NewJSONHandler(multiWriter, opts)

	return slog.New(handler)
}
