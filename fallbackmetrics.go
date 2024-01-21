package circuitotel

import (
	"context"
	"time"

	"github.com/cep21/circuit/v4"
)

type fallbackMetrics struct {
	wrappedSingletons
}

func (f *fallbackMetrics) Success(ctx context.Context, _ time.Time, duration time.Duration) {
	f.wrappedSingletons.histogramRecord(ctx, "fallback.success", f.m.fallbackSuccess, duration)
}

func (f *fallbackMetrics) ErrFailure(ctx context.Context, _ time.Time, duration time.Duration) {
	f.wrappedSingletons.histogramRecord(ctx, "fallback.error", f.m.fallbackErrFail, duration)
}

func (f *fallbackMetrics) ErrConcurrencyLimitReject(ctx context.Context, _ time.Time) {
	f.wrappedSingletons.counterIncr(ctx, "fallback.concurrency_limit", f.m.fallbackErrConcurrencyLimitReject)
}

var _ circuit.FallbackMetrics = &fallbackMetrics{}
