package methods

import "github.com/prometheus/prometheus/pkg/rulefmt"

type AlertErrorOptions struct {
	ServiceName        string
	AvailabilityTarget float64

	// important for simple algorithm
	AlertWindow string
	AlertWait   string
}

type AlertLatencyOptions struct {
	ServiceName string
	Targets     []LatencyTarget

	// important for simple algorithm
	AlertWindow string
	AlertWait   string
}

type AlertMethod interface {
	AlertForError(*AlertErrorOptions) ([]rulefmt.Rule, error)
	AlertForLatency(*AlertLatencyOptions) ([]rulefmt.Rule, error)
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
