package requestctx

import (
	"context"
	"strings"
	"testing"
	"time"
)

const (
	requestIDForTest = "req-789"
	traceIDForTest   = "trace-123"
	userIDForTest    = "user-001"
)

func TestWithInfoStoresRequestMetadata(t *testing.T) {
	start := time.Now()
	ctx := WithInfo(context.Background(), Info{
		TraceID:   traceIDForTest,
		SpanID:    "span-456",
		RequestID: requestIDForTest,
		UserID:    userIDForTest,
		StartTime: start,
	})

	info, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext() ok = false, want true")
	}
	if info.TraceID != traceIDForTest {
		t.Fatalf("TraceID = %q, want %s", info.TraceID, traceIDForTest)
	}
	if info.SpanID != "span-456" {
		t.Fatalf("SpanID = %q, want span-456", info.SpanID)
	}
	if info.RequestID != requestIDForTest {
		t.Fatalf("RequestID = %q, want %s", info.RequestID, requestIDForTest)
	}
	if info.UserID != userIDForTest {
		t.Fatalf("UserID = %q, want %s", info.UserID, userIDForTest)
	}
	if !info.StartTime.Equal(start) {
		t.Fatalf("StartTime = %v, want %v", info.StartTime, start)
	}
}

func TestWithInfoDefaultsStartTime(t *testing.T) {
	before := time.Now()
	ctx := WithInfo(context.Background(), Info{})
	after := time.Now()

	got := GetStartTime(ctx)
	if got.Before(before) || got.After(after) {
		t.Fatalf("StartTime = %v, want between %v and %v", got, before, after)
	}
}

func TestFromContextWithoutMetadata(t *testing.T) {
	info, ok := FromContext(context.Background())
	if ok {
		t.Fatal("FromContext() ok = true, want false")
	}
	if info != (Info{}) {
		t.Fatalf("Info = %+v, want zero value", info)
	}
}

func TestWithTraceInfoSupportsGetters(t *testing.T) {
	ctx := WithTraceInfo(context.Background(), traceIDForTest, "span-456", requestIDForTest)

	if got := GetTraceID(ctx); got != traceIDForTest {
		t.Fatalf("GetTraceID() = %q, want %s", got, traceIDForTest)
	}
	if got := GetRequestID(ctx); got != requestIDForTest {
		t.Fatalf("GetRequestID() = %q, want %s", got, requestIDForTest)
	}
	if GetStartTime(ctx).IsZero() {
		t.Fatal("GetStartTime() is zero")
	}
}

func TestWithTraceSpanPreservesExistingMetadata(t *testing.T) {
	ctx := WithInfo(context.Background(), Info{
		RequestID: requestIDForTest,
		UserID:    userIDForTest,
	})

	ctx = WithTraceSpan(ctx, traceIDForTest, "span-456")

	info, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext() ok = false, want true")
	}
	if info.RequestID != requestIDForTest {
		t.Fatalf("RequestID = %q, want %s", info.RequestID, requestIDForTest)
	}
	if info.UserID != userIDForTest {
		t.Fatalf("UserID = %q, want %s", info.UserID, userIDForTest)
	}
	if info.TraceID != traceIDForTest {
		t.Fatalf("TraceID = %q, want %s", info.TraceID, traceIDForTest)
	}
	if info.SpanID != "span-456" {
		t.Fatalf("SpanID = %q, want span-456", info.SpanID)
	}
}

func TestCleanMetadataID(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "", want: ""},
		{name: "visible ascii", value: "req-123", want: "req-123"},
		{name: "space is rejected", value: "req 123", want: ""},
		{name: "non ascii is rejected", value: "请求", want: ""},
		{name: "overlong is rejected", value: strings.Repeat("a", 129), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanMetadataID(tt.value); got != tt.want {
				t.Fatalf("CleanMetadataID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWithUserIDAddsAuthenticatedIdentity(t *testing.T) {
	start := time.Now()
	ctx := WithInfo(context.Background(), Info{
		TraceID:   traceIDForTest,
		RequestID: "req-789",
		StartTime: start,
	})

	ctx = WithUserID(ctx, "user-001")

	info, ok := FromContext(ctx)
	if !ok {
		t.Fatal("FromContext() ok = false, want true")
	}
	if info.TraceID != traceIDForTest {
		t.Fatalf("TraceID = %q, want %s", info.TraceID, traceIDForTest)
	}
	if info.UserID != "user-001" {
		t.Fatalf("UserID = %q, want user-001", info.UserID)
	}
	if !info.StartTime.Equal(start) {
		t.Fatalf("StartTime = %v, want %v", info.StartTime, start)
	}
}

func TestGettersHandleNilContext(t *testing.T) {
	var ctx context.Context
	if got := GetTraceID(ctx); got != "" {
		t.Fatalf("GetTraceID(nil) = %q, want empty", got)
	}
	if got := GetRequestID(ctx); got != "" {
		t.Fatalf("GetRequestID(nil) = %q, want empty", got)
	}
	if got := GetUserID(ctx); got != "" {
		t.Fatalf("GetUserID(nil) = %q, want empty", got)
	}
	if got := GetStartTime(ctx); !got.IsZero() {
		t.Fatalf("GetStartTime(nil) = %v, want zero", got)
	}
}
