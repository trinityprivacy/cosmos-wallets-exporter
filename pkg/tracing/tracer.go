package tracing

import (
	"context"
	"encoding/base64"
	"fmt"
	configPkg "main/pkg/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func getExporter(config configPkg.TracingConfig) (tracesdk.SpanExporter, error) {
	if config.Enabled.Bool {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(config.OpenTelemetryHTTPHost),
		}

		if config.OpenTelemetryHTTPInsecure.Bool {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		if config.OpenTelemetryHTTPUser != "" && config.OpenTelemetryHTTPPassword != "" {
			auth := config.OpenTelemetryHTTPUser + ":" + config.OpenTelemetryHTTPPassword
			token := base64.StdEncoding.EncodeToString([]byte(auth))
			opts = append(opts, otlptracehttp.WithHeaders(map[string]string{
				"Authorization": "Basic " + token,
			}))
		}

		return otlptracehttp.New(
			context.Background(),
			opts...,
		)
	}

	return NewNoopExporter(), nil
}

func InitTracer(config configPkg.TracingConfig, version string) (trace.Tracer, error) {
	exporter, err := getExporter(config)
	if err != nil {
		return nil, fmt.Errorf("error creating exporter: %w", err)
	}

	tp, err := NewTraceProvider(exporter, version)
	if err != nil {
		return nil, fmt.Errorf("error initializizng provider: %w", err)
	}

	otel.SetTracerProvider(tp)

	return tp.Tracer("main"), nil
}
