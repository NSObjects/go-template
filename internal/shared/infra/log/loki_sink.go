package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// LokiSink Loki输出
type LokiSink struct {
	client  *http.Client
	url     string
	labels  map[string]string
	timeout time.Duration
}

// LokiSinkConfig Loki输出配置
type LokiSinkConfig struct {
	URL     string            `json:"url" yaml:"url" toml:"url"`
	Labels  map[string]string `json:"labels" yaml:"labels" toml:"labels"`
	Timeout time.Duration     `json:"timeout" yaml:"timeout" toml:"timeout"`
}

func NewLokiSink(cfg LokiSinkConfig) *LokiSink {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	labels := cfg.Labels
	if labels == nil {
		labels = map[string]string{
			"service": "echo-admin",
		}
	}

	return &LokiSink{
		client:  &http.Client{Timeout: timeout},
		url:     cfg.URL,
		labels:  labels,
		timeout: timeout,
	}
}

func (l *LokiSink) Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error {
	// 构建Loki格式的日志条目
	entry := map[string]interface{}{
		"timestamp": time.Now().UnixNano(),
		"level":     level.String(),
		"message":   msg,
	}

	for _, attr := range attrs {
		entry[attr.Key] = attr.Value.Any()
	}

	// 构建Loki push API请求
	lokiEntry := map[string]interface{}{
		"stream": l.labels,
		"values": [][]string{
			{fmt.Sprintf("%d", time.Now().UnixNano()), fmt.Sprintf("%v", entry)},
		},
	}

	payload := map[string]interface{}{
		"streams": []interface{}{lokiEntry},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.url+"/loki/api/v1/push", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("loki request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (l *LokiSink) Close() error {
	return nil
}
