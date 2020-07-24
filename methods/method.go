package methods

import (
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type Window struct {
	Duration     model.Duration `yaml:"duration"`
	Consumption  float64        `yaml:"consumption"`
	Notification string         `yaml:"notification"`
}

type AlertErrorOptions struct {
	ServiceName        string
	AvailabilityTarget float64
	SLOWindow          time.Duration

	Windows     []Window
	ShortWindow bool

	// important for simple algorithm
	AlertWindow string
	AlertWait   string
}

type AlertLatencyOptions struct {
	ServiceName string
	Targets     []LatencyTarget
	SLOWindow   time.Duration

	Windows     []Window
	ShortWindow bool

	// important for simple algorithm
	AlertWindow string
	AlertWait   string
}

type AlertMethod interface {
	AlertForError(*AlertErrorOptions) ([]rulefmt.Rule, error)
	AlertForLatency(*AlertLatencyOptions) ([]rulefmt.Rule, error)
}

// Severities list of available severities: page and ticket
var Severities = []string{"page", "ticket"}

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
