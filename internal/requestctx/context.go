// Package requestctx carries request-scoped metadata without depending on an
// HTTP framework.
package requestctx

import (
	"context"
	"time"
)

type contextKey struct{}

// Info is metadata extracted at a delivery boundary and passed through
// request-scoped work.
type Info struct {
	TraceID   string
	SpanID    string
	RequestID string
	UserID    string
	StartTime time.Time
}

// WithInfo returns a child context carrying request metadata.
func WithInfo(ctx context.Context, info Info) context.Context {
	if info.StartTime.IsZero() {
		info.StartTime = time.Now()
	}
	return context.WithValue(ctx, contextKey{}, info)
}

// WithTraceInfo returns a child context carrying common trace metadata.
func WithTraceInfo(ctx context.Context, traceID, spanID, requestID, userID string) context.Context {
	return WithInfo(ctx, Info{
		TraceID:   traceID,
		SpanID:    spanID,
		RequestID: requestID,
		UserID:    userID,
	})
}

// FromContext returns request metadata from ctx.
func FromContext(ctx context.Context) (Info, bool) {
	if ctx == nil {
		return Info{}, false
	}
	info, ok := ctx.Value(contextKey{}).(Info)
	return info, ok
}

// GetTraceID returns the trace ID stored in ctx.
func GetTraceID(ctx context.Context) string {
	info, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	return info.TraceID
}

// GetRequestID returns the request ID stored in ctx.
func GetRequestID(ctx context.Context) string {
	info, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	return info.RequestID
}

// GetUserID returns the user ID stored in ctx.
func GetUserID(ctx context.Context) string {
	info, ok := FromContext(ctx)
	if !ok {
		return ""
	}
	return info.UserID
}

// GetStartTime returns the request start time stored in ctx.
func GetStartTime(ctx context.Context) time.Time {
	info, ok := FromContext(ctx)
	if !ok {
		return time.Time{}
	}
	return info.StartTime
}
