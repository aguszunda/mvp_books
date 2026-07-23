package metrics

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	HTTPRequestCount     metric.Int64Counter
	HTTPRequestDuration  metric.Float64Histogram
	HTTPRequestsInFlight metric.Int64UpDownCounter
)

func Init(meter metric.Meter) error {
	var err error

	HTTPRequestCount, err = meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	HTTPRequestDuration, err = meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return err
	}

	HTTPRequestsInFlight, err = meter.Int64UpDownCounter(
		"http.server.request.inflight",
		metric.WithDescription("Number of HTTP requests currently in progress"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	return nil
}

func RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration float64) {
	attrs := metric.WithAttributes(
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		semconv.HTTPResponseStatusCode(statusCode),
	)

	HTTPRequestCount.Add(ctx, 1, attrs)
	HTTPRequestDuration.Record(ctx, duration, attrs)
}

func IncInFlight(ctx context.Context) {
	HTTPRequestsInFlight.Add(ctx, 1)
}

func DecInFlight(ctx context.Context) {
	HTTPRequestsInFlight.Add(ctx, -1)
}
