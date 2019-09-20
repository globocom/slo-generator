package methods

import "github.com/prometheus/prometheus/pkg/rulefmt"

type AlertMethod interface {
	AlertForError(serviceName string, availabilityTarget float64) []rulefmt.Rule
	AlertForLatency(serviceName string, targets []LatencyTarget) []rulefmt.Rule
}

var methods = map[string]AlertMethod{}

func register(method AlertMethod, name string) AlertMethod {
	methods[name] = method
	return method
}

func Get(name string) AlertMethod {
	return methods[name]
}

type LatencyTarget struct {
	LE     string  `yaml:"le"`
	Target float64 `yaml:"target"`
}
