/*
 * OpenTelemetry 初始化模块
 * 提供分布式追踪功能
 *
 * Created by lintao on 2025/01/27
 * Copyright © 2020-2024 LINTAO. All rights reserved.
 */

package otel

import (
	"context"
	"fmt"
	"time"

	"github.com/NSObjects/go-kit/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
)

// TracerProvider 封装 OpenTelemetry TracerProvider
type TracerProvider struct {
	tp *sdktrace.TracerProvider
}

// Shutdown 关闭 TracerProvider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp.tp != nil {
		return tp.tp.Shutdown(ctx)
	}
	return nil
}

// NewTracerProvider 创建并初始化 TracerProvider
func NewTracerProvider(lc fx.Lifecycle, cfg config.OtelConfig) (*TracerProvider, error) {
	if !cfg.Enabled {
		// 如果未启用，返回空的 TracerProvider
		return &TracerProvider{tp: nil}, nil
	}

	ctx := context.Background()

	// 设置服务名称
	serviceName := cfg.ServiceName
	if serviceName == "" {
		serviceName = "go-template"
	}

	// 设置环境
	environment := cfg.Environment
	if environment == "" {
		environment = "dev"
	}

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentName(environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建导出器
	var exporter sdktrace.SpanExporter
	switch cfg.ExporterType {
	case "otlp":
		endpoint := cfg.OTLPEndpoint
		if endpoint == "" {
			endpoint = "localhost:4317"
		}
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(), // 开发环境使用，生产环境应使用 TLS
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
	default:
		// 默认使用 stdout，便于开发调试
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
	}

	// 设置采样率
	samplingRatio := cfg.SamplingRatio
	if samplingRatio <= 0 {
		samplingRatio = 1.0 // 默认 100% 采样
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(samplingRatio)),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracerProvider := &TracerProvider{tp: tp}

	// 注册生命周期钩子
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return tracerProvider.Shutdown(shutdownCtx)
		},
	})

	return tracerProvider, nil
}

// Module Fx 模块
var Module = fx.Options(
	fx.Provide(NewTracerProvider),
)
