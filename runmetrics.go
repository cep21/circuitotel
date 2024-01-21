package circuitotel

import (
	"context"
	"time"

	"github.com/cep21/circuit/v4"
)

type runMetrics struct {
	wrappedSingletons
}

func (r *runMetrics) Success(ctx context.Context, _ time.Time, duration time.Duration) {
	r.wrappedSingletons.histogramRecord(ctx, "run.success", r.m.runSuccess, duration)
}

func (r *runMetrics) ErrFailure(ctx context.Context, _ time.Time, duration time.Duration) {
	r.wrappedSingletons.histogramRecord(ctx, "run.failure", r.m.runErrFail, duration)
}

func (r *runMetrics) ErrTimeout(ctx context.Context, _ time.Time, duration time.Duration) {
	r.wrappedSingletons.histogramRecord(ctx, "run.timeout", r.m.errTimeout, duration)
}

func (r *runMetrics) ErrBadRequest(ctx context.Context, _ time.Time, duration time.Duration) {
	r.wrappedSingletons.histogramRecord(ctx, "run.bad_request", r.m.errBadRequest, duration)
}

func (r *runMetrics) ErrInterrupt(ctx context.Context, _ time.Time, duration time.Duration) {
	r.wrappedSingletons.histogramRecord(ctx, "run.interrupt", r.m.errInterrupt, duration)
}

func (r *runMetrics) ErrConcurrencyLimitReject(ctx context.Context, _ time.Time) {
	r.wrappedSingletons.counterIncr(ctx, "run.concurrency_limit", r.m.runErrConcurrencyLimitReject)
}

func (r *runMetrics) ErrShortCircuit(ctx context.Context, _ time.Time) {
	r.wrappedSingletons.counterIncr(ctx, "run.short_circuit", r.m.errShortCircuit)
}

var _ circuit.RunMetrics = &runMetrics{}
