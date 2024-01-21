package circuitotel_test

import (
	"context"

	"github.com/cep21/circuit/v4"
	"github.com/cep21/circuitotel"
)

func ExampleFactory_simple() {
	// For many, the default factory is sufficient
	var f circuitotel.Factory
	manager := circuit.Manager{
		// Pass the factory constructor to your manager
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			f.CommandPropertiesConstructor,
		},
	}
	_ = manager.MustCreateCircuit("test-circuit")
}

func ExampleFactory() {
	// Make a factory.  With defaults, it uses the
	// global meter provider.
	factory := circuitotel.Factory{
		// You can pass your own provider if you wish
		// MeterProvider: nil
	}
	manager := circuit.Manager{
		// Pass the factory constructor to your manager
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			factory.CommandPropertiesConstructor,
		},
	}
	// Now make and use circuits like normal
	c := manager.MustCreateCircuit("test-circuit")
	_ = c.Execute(context.Background(), func(_ context.Context) error {
		return nil
	}, nil)
}
