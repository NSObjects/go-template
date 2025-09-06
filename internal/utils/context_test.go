/*
 * Context工具函数测试
 * 测试链路追踪上下文相关功能
 */

package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestExtractTraceContext(t *testing.T) {
	tests := []struct {
		name          string
		headers       map[string]string
		userID        interface{}
		expectedTrace *TraceContext
	}{
		{
			name: "complete trace info",
			headers: map[string]string{
				"X-Request-ID": "req-123",
				"X-Trace-ID":   "trace-456",
				"X-Span-ID":    "span-789",
			},
			userID: "user-001",
			expectedTrace: &TraceContext{
				TraceID:   "trace-456",
				SpanID:    "span-789",
				RequestID: "req-123",
				UserID:    "user-001",
			},
		},
		{
			name: "minimal trace info",
			headers: map[string]string{
				"X-Request-ID": "req-456",
			},
			userID: nil,
			expectedTrace: &TraceContext{
				TraceID:   "",
				SpanID:    "",
				RequestID: "req-456",
				UserID:    "",
			},
		},
		{
			name:    "no trace info",
			headers: map[string]string{},
			userID:  nil,
			expectedTrace: &TraceContext{
				TraceID:   "",
				SpanID:    "",
				RequestID: "",
				UserID:    "",
			},
		},
		{
			name: "invalid user ID type",
			headers: map[string]string{
				"X-Request-ID": "req-789",
			},
			userID: 123, // 不是string类型
			expectedTrace: &TraceContext{
				TraceID:   "",
				SpanID:    "",
				RequestID: "req-789",
				UserID:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Echo实例
			e := echo.New()

			// 创建测试请求
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 设置用户ID
			if tt.userID != nil {
				c.Set("user_id", tt.userID)
			}

			// 提取链路追踪信息
			trace := ExtractTraceContext(c)

			// 验证结果
			assert.Equal(t, tt.expectedTrace.TraceID, trace.TraceID)
			assert.Equal(t, tt.expectedTrace.SpanID, trace.SpanID)
			assert.Equal(t, tt.expectedTrace.RequestID, trace.RequestID)
			assert.Equal(t, tt.expectedTrace.UserID, trace.UserID)
			assert.NotZero(t, trace.StartTime)
		})
	}
}

func TestBuildContext(t *testing.T) {
	// 创建Echo实例
	e := echo.New()

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "req-123")
	req.Header.Set("X-Trace-ID", "trace-456")
	req.Header.Set("X-Span-ID", "span-789")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-001")

	// 构建上下文
	ctx := BuildContext(c)

	// 验证上下文中的值
	assert.Equal(t, "trace-456", GetTraceID(ctx))
	assert.Equal(t, "span-789", ctx.Value("span_id"))
	assert.Equal(t, "req-123", GetRequestID(ctx))
	assert.Equal(t, "user-001", GetUserID(ctx))
	assert.NotZero(t, GetStartTime(ctx))
}

func TestGetTraceID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with trace ID",
			ctx:      context.WithValue(context.Background(), "trace_id", "trace-123"),
			expected: "trace-123",
		},
		{
			name:     "without trace ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "invalid trace ID type",
			ctx:      context.WithValue(context.Background(), "trace_id", 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTraceID(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with request ID",
			ctx:      context.WithValue(context.Background(), "request_id", "req-123"),
			expected: "req-123",
		},
		{
			name:     "without request ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "invalid request ID type",
			ctx:      context.WithValue(context.Background(), "request_id", 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRequestID(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "with user ID",
			ctx:      context.WithValue(context.Background(), "user_id", "user-123"),
			expected: "user-123",
		},
		{
			name:     "without user ID",
			ctx:      context.Background(),
			expected: "",
		},
		{
			name:     "invalid user ID type",
			ctx:      context.WithValue(context.Background(), "user_id", 123),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserID(tt.ctx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStartTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		ctx      context.Context
		expected time.Time
	}{
		{
			name:     "with start time",
			ctx:      context.WithValue(context.Background(), "start_time", now),
			expected: now,
		},
		{
			name:     "without start time",
			ctx:      context.Background(),
			expected: time.Now(), // 应该返回当前时间
		},
		{
			name:     "invalid start time type",
			ctx:      context.WithValue(context.Background(), "start_time", "invalid"),
			expected: time.Now(), // 应该返回当前时间
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStartTime(tt.ctx)
			if tt.name == "with start time" {
				assert.Equal(t, tt.expected, result)
			} else {
				// 对于没有start_time或类型错误的情况，应该返回当前时间
				assert.True(t, result.After(now.Add(-time.Second)))
				assert.True(t, result.Before(now.Add(time.Second)))
			}
		})
	}
}

func TestWithTraceInfo(t *testing.T) {
	ctx := context.Background()
	traceID := "trace-123"
	spanID := "span-456"
	requestID := "req-789"
	userID := "user-001"

	// 添加链路追踪信息
	newCtx := WithTraceInfo(ctx, traceID, spanID, requestID, userID)

	// 验证添加的信息
	assert.Equal(t, traceID, GetTraceID(newCtx))
	assert.Equal(t, spanID, newCtx.Value("span_id"))
	assert.Equal(t, requestID, GetRequestID(newCtx))
	assert.Equal(t, userID, GetUserID(newCtx))
	assert.NotZero(t, GetStartTime(newCtx))
}

func TestTraceContext_Fields(t *testing.T) {
	now := time.Now()
	trace := &TraceContext{
		TraceID:   "trace-123",
		SpanID:    "span-456",
		RequestID: "req-789",
		UserID:    "user-001",
		StartTime: now,
	}

	// 验证字段
	assert.Equal(t, "trace-123", trace.TraceID)
	assert.Equal(t, "span-456", trace.SpanID)
	assert.Equal(t, "req-789", trace.RequestID)
	assert.Equal(t, "user-001", trace.UserID)
	assert.Equal(t, now, trace.StartTime)
}

func TestExtractTraceContext_ResponseHeader(t *testing.T) {
	// 创建Echo实例
	e := echo.New()

	// 创建测试请求
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 在响应头中设置RequestID
	c.Response().Header().Set("X-Request-ID", "resp-req-123")

	// 提取链路追踪信息
	trace := ExtractTraceContext(c)

	// 验证从响应头中提取的RequestID
	assert.Equal(t, "resp-req-123", trace.RequestID)
}
