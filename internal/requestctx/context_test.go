package requestctx

import (
	"context"
	"testing"
	"time"
)

const traceIDForTest = "trace-123"

func TestWithInfoStoresRequestMetadata(t *testing.T) {
	start := time.Now()
	ctx := WithInfo(context.Background(), Info{
		TraceID:   traceIDForTest,
		SpanID:    "span-456",
		RequestID: "req-789",
		UserID:    "user-001",
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
	if info.RequestID != "req-789" {
		t.Fatalf("RequestID = %q, want req-789", info.RequestID)
	}
	if info.UserID != "user-001" {
		t.Fatalf("UserID = %q, want user-001", info.UserID)
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
	ctx := WithTraceInfo(context.Background(), traceIDForTest, "span-456", "req-789")

	if got := GetTraceID(ctx); got != traceIDForTest {
		t.Fatalf("GetTraceID() = %q, want %s", got, traceIDForTest)
	}
	if got := GetRequestID(ctx); got != "req-789" {
		t.Fatalf("GetRequestID() = %q, want req-789", got)
	}
	if GetStartTime(ctx).IsZero() {
		t.Fatal("GetStartTime() is zero")
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
