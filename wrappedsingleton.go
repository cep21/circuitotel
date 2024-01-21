package circuitotel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type wrappedSingletons struct {
	m        *metricSingletons
	attrs    attribute.Set
	attrList []attribute.KeyValue
}

func (w *wrappedSingletons) histogramRecord(ctx context.Context, eventName string, histogram metric.Float64Histogram, duration time.Duration) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("circuit."+eventName, trace.WithAttributes(w.attrList...))
	histogram.Record(ctx, float64(duration.Nanoseconds())/float64(time.Millisecond), metric.WithAttributeSet(w.attrs))
}

func (w *wrappedSingletons) counterIncr(ctx context.Context, eventName string, counter metric.Int64Counter) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("circuit."+eventName, trace.WithAttributes(w.attrList...))
	counter.Add(ctx, 1, metric.WithAttributeSet(w.attrs))
}
