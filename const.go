package circuitotel

import "go.opentelemetry.io/otel/attribute"

const ScopeName = "github.com/cep21/circuitotel"
const attrName = attribute.Key("circuit.name")

func Version() string {
	return "0.0.1" // TODO: Auto update on release
}
