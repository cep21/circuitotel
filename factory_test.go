package circuitotel

import (
	"context"
	"testing"

	"github.com/cep21/circuit/v4"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestFactory(t *testing.T) {
	sr := tracetest.NewSpanRecorder()
	tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))

	reader := metric.NewManualReader()
	meterProvider := metric.NewMeterProvider(metric.WithReader(reader))

	f := &Factory{
		MeterProvider: meterProvider,
	}
	cm := circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			f.CommandPropertiesConstructor,
		},
	}
	c := cm.MustCreateCircuit("test-circuit")
	val := 0
	ctx := context.Background()
	ctx, sp := tp.Tracer("test").Start(ctx, "test")
	require.NoError(t, c.Execute(ctx, func(ctx context.Context) error {
		val++
		return nil
	}, nil))
	sp.End()
	require.Equal(t, 1, val)
	require.Equal(t, 1, len(sr.Ended()))
	endedSpan := sr.Ended()[0]
	require.Equal(t, "test", endedSpan.InstrumentationScope().Name)
	require.Equal(t, "test", endedSpan.Name())
	require.Len(t, endedSpan.Events(), 1)
	require.Equal(t, "circuit.run.success", endedSpan.Events()[0].Name)
	require.NoError(t, sr.Shutdown(ctx))
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Equal(t, ScopeName, rm.ScopeMetrics[0].Scope.Name)
	x := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64]).DataPoints[0].Count
	circName, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Histogram[float64]).DataPoints[0].Attributes.Value(attrName)
	require.True(t, ok)
	require.Equal(t, "test-circuit", circName.AsString())
	require.Equal(t, uint64(1), x)
	require.NoError(t, reader.Shutdown(ctx))
}
