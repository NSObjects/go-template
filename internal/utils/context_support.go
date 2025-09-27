/*
 * Context工具函数支持
 * 提供从 echo.Context 提取链路追踪信息并构建标准 context.Context 的能力。
 */
package utils

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

// TraceContext 链路追踪上下文信息
type TraceContext struct {
	TraceID   string
	SpanID    string
	RequestID string
	UserID    string
	StartTime time.Time
}

// ExtractTraceContext 从 echo.Context 中提取链路追踪信息
func ExtractTraceContext(c echo.Context) *TraceContext {
	tc := &TraceContext{
		StartTime: time.Now(),
	}

	// 提取请求 ID
	if requestID := c.Request().Header.Get("X-Request-ID"); requestID != "" {
		tc.RequestID = requestID
	} else if requestID := c.Response().Header().Get("X-Request-ID"); requestID != "" {
		tc.RequestID = requestID
	}

	// 提取 TraceID 和 SpanID (OpenTelemetry 格式)
	if traceID := c.Request().Header.Get("X-Trace-ID"); traceID != "" {
		tc.TraceID = traceID
	}
	if spanID := c.Request().Header.Get("X-Span-ID"); spanID != "" {
		tc.SpanID = spanID
	}

	// 提取用户 ID (如果已认证)
	if userID := c.Get("user_id"); userID != nil {
		if uid, ok := userID.(string); ok {
			tc.UserID = uid
		}
	}

	return tc
}

// BuildContext 构造包含链路追踪信息的标准 context.Context
func BuildContext(c echo.Context) context.Context {
	ctx := c.Request().Context()
	tc := ExtractTraceContext(c)

	// 将链路追踪信息添加到 context 中
	ctx = context.WithValue(ctx, "trace_id", tc.TraceID)
	ctx = context.WithValue(ctx, "span_id", tc.SpanID)
	ctx = context.WithValue(ctx, "request_id", tc.RequestID)
	ctx = context.WithValue(ctx, "user_id", tc.UserID)
	ctx = context.WithValue(ctx, "start_time", tc.StartTime)

	return ctx
}

// GetTraceID 从 context 中获取 TraceID
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// GetRequestID 从 context 中获取 RequestID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}

// GetUserID 从 context 中获取 UserID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// GetStartTime 从 context 中获取请求开始时间
func GetStartTime(ctx context.Context) time.Time {
	if startTime, ok := ctx.Value("start_time").(time.Time); ok {
		return startTime
	}
	return time.Now()
}

// WithTraceInfo 为 context 添加链路追踪信息
func WithTraceInfo(ctx context.Context, traceID, spanID, requestID, userID string) context.Context {
	ctx = context.WithValue(ctx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "span_id", spanID)
	ctx = context.WithValue(ctx, "request_id", requestID)
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "start_time", time.Now())
	return ctx
}
