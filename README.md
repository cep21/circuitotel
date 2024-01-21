# circuitotel
An open telemetry implementation of [circuit](https://github.com/cep21/circuit) library metric collector

## Usage

For many users, the defaults will work just fine.  You can also pass
a custom metric provider to your factory.

```go
	// For many, the default factory is sufficient
	var f circuitotel.Factory
	manager := circuit.Manager{
		// Pass the factory constructor to your manager
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			f.CommandPropertiesConstructor,
		},
	}
	_ = manager.MustCreateCircuit("test-circuit")
```

## Traces

We do not start traces since it's expected the downstream itself will trace.
Instead, we pass circuit events to the passed in span.

## Metrics

We gather all circuit exposed metrics as histograms or counters