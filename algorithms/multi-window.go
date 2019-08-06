package algoritms

import (
	"fmt"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type MultiWindowAlgorithm struct{}

func (*MultiWindowAlgorithm) AlertForError(serviceName string, availabilityTarget float64, annotations map[string]string) []rulefmt.Rule {
	rules := []rulefmt.Rule{
		{
			Alert: "slo:" + serviceName + ".errors.page",
			Expr: multiBurnRate(MultiRateOpts{
				Metric:   "slo:service_errors_total",
				Labels:   labels.New(labels.Label{"service", serviceName}),
				Operator: ">",
				Value:    (1 - availabilityTarget/100),
				Kind:     "page",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "page",
			},
		},
		{
			Alert: "slo:" + serviceName + ".errors.ticket",
			Expr: multiBurnRate(MultiRateOpts{
				Metric:   "slo:service_errors_total",
				Labels:   labels.New(labels.Label{"service", serviceName}),
				Operator: ">",
				Value:    (1 - availabilityTarget/100),
				Kind:     "ticket",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "ticket",
			},
		},
	}
	// rulefmt.Rule
	return rules
}

type MultiRateOpts struct {
	Metric   string
	Labels   labels.Labels
	Operator string
	Value    float64
	Kind     string // page or ticket
}

type MultiRateWindow [2]struct {
	Multiplier  float64
	LongWindow  string
	ShortWindow string
}

var multiRateWindows = map[string]MultiRateWindow{
	"page": {
		{
			Multiplier:  14.4,
			LongWindow:  "1h",
			ShortWindow: "5m",
		},
		{
			Multiplier:  6,
			LongWindow:  "6h",
			ShortWindow: "30m",
		},
	},
	"ticket": {
		{
			Multiplier:  3,
			LongWindow:  "1d",
			ShortWindow: "2h",
		},
		{
			Multiplier:  1,
			LongWindow:  "3d",
			ShortWindow: "6h",
		},
	},
}

func multiBurnRate(opts MultiRateOpts) string {
	multiRateWindow := multiRateWindows[opts.Kind]

	result := ""
	result += fmt.Sprintf(`(%s:ratio_rate_%s%s %s (%g * %.3g) and `, opts.Metric, multiRateWindow[0].LongWindow, opts.Labels.String(), opts.Operator, multiRateWindow[0].Multiplier, opts.Value)
	result += fmt.Sprintf(`%s:ratio_rate_%s%s %s (%g * %.3g))`, opts.Metric, multiRateWindow[0].ShortWindow, opts.Labels.String(), opts.Operator, multiRateWindow[0].Multiplier, opts.Value)

	result += " or "

	result += fmt.Sprintf(`(%s:ratio_rate_%s%s %s (%g * %.3g) and `, opts.Metric, multiRateWindow[1].LongWindow, opts.Labels.String(), opts.Operator, multiRateWindow[1].Multiplier, opts.Value)
	result += fmt.Sprintf(`%s:ratio_rate_%s%s %s (%g * %.3g))`, opts.Metric, multiRateWindow[1].ShortWindow, opts.Labels.String(), opts.Operator, multiRateWindow[1].Multiplier, opts.Value)

	return result
}

var _ = register(&MultiWindowAlgorithm{}, "multi-window")
