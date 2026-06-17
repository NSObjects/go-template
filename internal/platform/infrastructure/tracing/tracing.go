// Package tracing owns the process-level OpenTelemetry runtime.
package tracing

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otelresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/NSObjects/go-template/internal/platform/configs"
	"github.com/NSObjects/go-template/internal/platform/infrastructure/resources"
)

// Runtime wraps an optional OpenTelemetry tracer provider and shutdown hook.
type Runtime struct {
	enabled         bool
	provider        *sdktrace.TracerProvider
	shutdown        func(context.Context) error
	shutdownTimeout time.Duration
}

// Start creates the configured tracing runtime. Disabled tracing requires no exporter.
func Start(ctx context.Context, cfg configs.TracingConfig, serviceName string) (*Runtime, error) {
	if ctx == nil {
		return nil, resources.NewCapabilityError(resources.CapabilityTracing, "start", errors.New("nil context"))
	}
	cfg = configs.Normalize(configs.Config{Tracing: cfg}).Tracing
	if !cfg.Enabled {
		return &Runtime{}, nil
	}
	if err := configs.Validate(configs.Config{Tracing: cfg}); err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityTracing, "configure", err)
	}
	if serviceName == "" {
		serviceName = configs.DefaultAppName
	}
	if cfg.ServiceName != "" {
		serviceName = cfg.ServiceName
	}

	exporter, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, resources.NewCapabilityError(resources.CapabilityTracing, "exporter", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(otelresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("telemetry.backend", "jaeger"),
		)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRatio)),
	)
	otel.SetTracerProvider(provider)

	return &Runtime{
		enabled:         true,
		provider:        provider,
		shutdown:        provider.Shutdown,
		shutdownTimeout: time.Duration(cfg.ShutdownTimeoutSeconds) * time.Second,
	}, nil
}

// Enabled reports whether tracing was enabled by configuration.
func (r *Runtime) Enabled() bool {
	return r != nil && r.enabled
}

// Check returns the current tracing capability status.
func (r *Runtime) Check(context.Context) resources.CapabilityStatus {
	if r == nil || !r.enabled {
		return resources.Disabled(resources.CapabilityTracing)
	}
	if r.provider == nil && r.shutdown == nil {
		return resources.Unavailable(resources.CapabilityTracing, errors.New("tracer provider is nil"))
	}
	return resources.Available(resources.CapabilityTracing, "provider installed")
}

// Shutdown flushes and releases the tracing provider.
func (r *Runtime) Shutdown(ctx context.Context) error {
	if r == nil || !r.enabled {
		return nil
	}
	if ctx == nil {
		return resources.NewCapabilityError(resources.CapabilityTracing, "shutdown", errors.New("nil context"))
	}
	if r.shutdownTimeout > 0 {
		shutdownCtx, cancel := context.WithTimeout(ctx, r.shutdownTimeout)
		defer cancel()
		ctx = shutdownCtx
	}
	if r.shutdown == nil {
		return resources.NewCapabilityError(resources.CapabilityTracing, "shutdown", errors.New("shutdown hook is nil"))
	}
	if err := r.shutdown(ctx); err != nil {
		return resources.NewCapabilityError(resources.CapabilityTracing, "shutdown", err)
	}
	otel.SetTracerProvider(noop.NewTracerProvider())
	return nil
}

func newExporter(ctx context.Context, cfg configs.TracingConfig) (*otlptrace.Exporter, error) {
	switch cfg.Protocol {
	case "http":
		options := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			options = append(options, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, options...)
	default:
		options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			options = append(options, otlptracegrpc.WithInsecure())
		}
		return otlptracegrpc.New(ctx, options...)
	}
}
