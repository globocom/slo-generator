package methods

import (
	"fmt"
	"strings"

	samples "github.com/globocom/slo-generator/samples"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type SimpleAlgorithm struct{}

func (*SimpleAlgorithm) AlertForError(opts *AlertErrorOptions) ([]rulefmt.Rule, error) {
	var (
		burnRate float64 = 1
		waitFor  model.Duration
	)

	if opts.BurnRate > 0 {
		burnRate = opts.BurnRate
	}

	if err := samples.ValidateSample(opts.AlertWindow); err != nil {
		return nil, err
	}
	if opts.AlertWait != "" {
		var err error
		waitFor, err = model.ParseDuration(opts.AlertWait)
		if err != nil {
			return nil, err
		}
	}

	ruleLabels := labels.New(labels.Label{Name: "service", Value: opts.ServiceName})
	errorLimit := 1 - opts.AvailabilityTarget/100
	rules := []rulefmt.Rule{
		{
			Alert:       "slo:" + opts.ServiceName + ".errors.page",
			Expr:        fmt.Sprintf("slo:service_errors_total:ratio_rate_%s%s > %.3g * %.3g", opts.AlertWindow, ruleLabels.String(), burnRate, errorLimit),
			For:         waitFor,
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "page",
				"signal":   "error",
			},
		},
	}

	return rules, nil
}

func (*SimpleAlgorithm) AlertForLatency(opts *AlertLatencyOptions) ([]rulefmt.Rule, error) {
	var waitFor model.Duration

	if err := samples.ValidateSample(opts.AlertWindow); err != nil {
		return nil, err
	}
	if opts.AlertWait != "" {
		var err error
		waitFor, err = model.ParseDuration(opts.AlertWait)
		if err != nil {
			return nil, err
		}
	}

	rules := []rulefmt.Rule{
		{
			Alert:       "slo:" + opts.ServiceName + ".latency.page",
			Expr:        simpleLatency(opts),
			For:         waitFor,
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "page",
				"signal":   "latency",
			},
		},
	}

	return rules, nil
}

func simpleLatency(opts *AlertLatencyOptions) string {
	var (
		conditions []string
		burnRate   float64 = 1
	)

	if opts.BurnRate > 0 {
		burnRate = opts.BurnRate
	}

	for _, target := range opts.Targets {
		value := 1 - ((100 - target.Target) * 0.01)

		lbs := labels.New(labels.Label{Name: "service", Value: opts.ServiceName}, labels.Label{Name: "le", Value: target.LE})
		condition := fmt.Sprintf(`slo:service_latency:ratio_rate_%s%s < %.3g * %.3g`, opts.AlertWindow, lbs.String(), burnRate, value)

		conditions = append(conditions, condition)
	}

	return strings.Join(conditions, " or ")
}

var _ = register(&SimpleAlgorithm{}, "simple")
