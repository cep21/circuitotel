package circuitotel

import (
	"context"
	"time"

	"github.com/cep21/circuit/v4"
)

type circuitMetrics struct {
	wrappedSingletons
}

func (c *circuitMetrics) Closed(ctx context.Context, _ time.Time) {
	c.wrappedSingletons.counterIncr(ctx, "closed", c.m.closed)
}

func (c *circuitMetrics) Opened(ctx context.Context, _ time.Time) {
	c.wrappedSingletons.counterIncr(ctx, "opened", c.m.open)
}

var _ circuit.Metrics = &circuitMetrics{}
