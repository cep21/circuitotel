package circuitotel

import (
	"fmt"
	"sync"

	"github.com/cep21/circuit/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

var _ circuit.CommandPropertiesConstructor = (&Factory{}).CommandPropertiesConstructor

type Factory struct {
	MeterProvider metric.MeterProvider

	mu               sync.Mutex
	metricSingletons *metricSingletons
}

func (f *Factory) getMeter() metric.Meter {
	mp := f.MeterProvider
	if mp == nil {
		mp = otel.GetMeterProvider()
	}
	return mp.Meter(ScopeName,
		// Note: May change these two attributes later: not 100% sure how they play into open
		//       telemetry
		metric.WithInstrumentationVersion(Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)
}

func (f *Factory) getMetricSingletons() (*metricSingletons, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.metricSingletons != nil {
		return f.metricSingletons, nil
	}
	m := f.getMeter()
	closed, err := m.Int64Counter("circuit.closed",
		metric.WithDescription("Counts the number of times a circuit is closed."),
		metric.WithUnit("{count}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.closed counter: %w", err)
	}
	open, err := m.Int64Counter("circuit.open",
		metric.WithDescription("Counts the number of times a circuit is open."),
		metric.WithUnit("{count}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.open counter: %w", err)
	}
	runSuccess, err := m.Float64Histogram("circuit.run.success",
		metric.WithDescription("Measures the duration of successful runs."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.success histogram: %w", err)
	}
	runErrFail, err := m.Float64Histogram("circuit.run.failure",
		metric.WithDescription("Measures the duration of failed runs."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.failure histogram: %w", err)
	}
	errTimeout, err := m.Float64Histogram("circuit.run.timeout",
		metric.WithDescription("Measures the duration of timed out runs."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.timeout histogram: %w", err)
	}
	errBadRequest, err := m.Float64Histogram("circuit.run.bad_request",
		metric.WithDescription("Measures the duration of bad request runs."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.bad_request histogram: %w", err)
	}
	errInterrupt, err := m.Float64Histogram("circuit.run.interrupt",
		metric.WithDescription("Measures the duration of interrupted runs."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.interrupt histogram: %w", err)
	}
	runErrConcurrencyLimitReject, err := m.Int64Counter("circuit.run.concurrency_limit",
		metric.WithDescription("Counts the number of times a run is rejected due to concurrency limit."),
		metric.WithUnit("{count}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.concurrency_limit counter: %w", err)
	}
	errShortCircuit, err := m.Int64Counter("circuit.run.short_circuit",
		metric.WithDescription("Counts the number of times a run is rejected due to short circuit."),
		metric.WithUnit("{count}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.run.short_circuit counter: %w", err)
	}
	fallbackSuccess, err := m.Float64Histogram("circuit.fallback.success",
		metric.WithDescription("Measures the duration of successful fallbacks."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.fallback.success histogram: %w", err)
	}
	fallbackErrFail, err := m.Float64Histogram("circuit.fallback.error",
		metric.WithDescription("Measures the duration of failed fallbacks."),
		metric.WithUnit("ms"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.fallback.error histogram: %w", err)
	}
	fallbackErrConcurrencyLimitReject, err := m.Int64Counter("circuit.fallback.concurrency_limit",
		metric.WithDescription("Counts the number of times a fallback is rejected due to concurrency limit."),
		metric.WithUnit("{count}"))
	if err != nil {
		return nil, fmt.Errorf("failed to create circuit.fallback.concurrency_limit counter: %w", err)
	}
	f.metricSingletons = &metricSingletons{
		closed:                            closed,
		open:                              open,
		runSuccess:                        runSuccess,
		runErrFail:                        runErrFail,
		errTimeout:                        errTimeout,
		errBadRequest:                     errBadRequest,
		errInterrupt:                      errInterrupt,
		runErrConcurrencyLimitReject:      runErrConcurrencyLimitReject,
		errShortCircuit:                   errShortCircuit,
		fallbackSuccess:                   fallbackSuccess,
		fallbackErrFail:                   fallbackErrFail,
		fallbackErrConcurrencyLimitReject: fallbackErrConcurrencyLimitReject,
	}
	return f.metricSingletons, nil
}

func (f *Factory) CommandPropertiesConstructor(name string) circuit.Config {
	cm, err := f.getMetricSingletons()
	if err != nil {
		otel.Handle(err)
		return circuit.Config{}
	}
	w := wrappedSingletons{
		m:     cm,
		attrs: attribute.NewSet(attrName.String(name)),
	}
	return circuit.Config{
		Metrics: circuit.MetricsCollectors{
			Circuit: []circuit.Metrics{
				&circuitMetrics{
					wrappedSingletons: w,
				},
			},
			Run: []circuit.RunMetrics{
				&runMetrics{
					wrappedSingletons: w,
				},
			},
			Fallback: []circuit.FallbackMetrics{
				&fallbackMetrics{
					wrappedSingletons: w,
				},
			},
		},
	}
}
