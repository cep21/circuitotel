package circuitotel

import "go.opentelemetry.io/otel/metric"

type metricSingletons struct {
	closed metric.Int64Counter
	open   metric.Int64Counter

	runSuccess                   metric.Float64Histogram
	runErrFail                   metric.Float64Histogram
	errTimeout                   metric.Float64Histogram
	errBadRequest                metric.Float64Histogram
	errInterrupt                 metric.Float64Histogram
	runErrConcurrencyLimitReject metric.Int64Counter
	errShortCircuit              metric.Int64Counter

	fallbackSuccess                   metric.Float64Histogram
	fallbackErrFail                   metric.Float64Histogram
	fallbackErrConcurrencyLimitReject metric.Int64Counter
}
