package algoritms

import "github.com/prometheus/prometheus/pkg/rulefmt"

type AlertAlgorithm interface {
	AlertForError(serviceName string, availabilityTarget float64, annotations map[string]string) []rulefmt.Rule
	AlertForLatency(serviceName string, targets []LatencyTarget, annotations map[string]string) []rulefmt.Rule
}

var algorithms = map[string]AlertAlgorithm{}

func register(algorithm AlertAlgorithm, name string) AlertAlgorithm {
	algorithms[name] = algorithm
	return algorithm
}

func Get(name string) AlertAlgorithm {
	return algorithms[name]
}

type LatencyTarget struct {
	LE     string  `yaml:"le"`
	Target float64 `yaml:"target"`
}
