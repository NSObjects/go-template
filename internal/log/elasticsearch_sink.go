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

// ElasticsearchSink Elasticsearch输出
type ElasticsearchSink struct {
	client  *http.Client
	url     string
	index   string
	timeout time.Duration
}

// ElasticsearchSinkConfig Elasticsearch输出配置
type ElasticsearchSinkConfig struct {
	URL     string        `json:"url" yaml:"url" toml:"url"`
	Index   string        `json:"index" yaml:"index" toml:"index"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" toml:"timeout"`
}

func NewElasticsearchSink(cfg ElasticsearchSinkConfig) *ElasticsearchSink {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	index := cfg.Index
	if index == "" {
		index = "echo-admin-logs"
	}

	return &ElasticsearchSink{
		client:  &http.Client{Timeout: timeout},
		url:     cfg.URL,
		index:   index,
		timeout: timeout,
	}
}

func (e *ElasticsearchSink) Write(ctx context.Context, level slog.Level, msg string, attrs []slog.Attr) error {
	entry := map[string]interface{}{
		"@timestamp": time.Now().Format(time.RFC3339),
		"level":      level.String(),
		"message":    msg,
	}

	for _, attr := range attrs {
		entry[attr.Key] = attr.Value.Any()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	// 构建ES bulk API请求
	bulkData := fmt.Sprintf(`{"index":{"_index":"%s"}}%s%s`, e.index, "\n", string(data))

	req, err := http.NewRequestWithContext(ctx, "POST", e.url+"/_bulk", bytes.NewBufferString(bulkData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("elasticsearch request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (e *ElasticsearchSink) Close() error {
	return nil
}
