package algoritms

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type MultiWindowAlgorithm struct{}

func (*MultiWindowAlgorithm) AlertForError(serviceName string, availabilityTarget float64, annotations map[string]string) []rulefmt.Rule {
	rules := []rulefmt.Rule{
		{
			Alert: "slo:" + serviceName + ".errors.page",
			Expr: multiBurnRate(MultiRateErrorOpts{
				Metric: "slo:service_errors_total",
				Labels: labels.New(labels.Label{"service", serviceName}),
				Value:  (1 - availabilityTarget/100),
				Kind:   "page",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "page",
			},
		},
		{
			Alert: "slo:" + serviceName + ".errors.ticket",
			Expr: multiBurnRate(MultiRateErrorOpts{
				Metric: "slo:service_errors_total",
				Labels: labels.New(labels.Label{"service", serviceName}),
				Value:  (1 - availabilityTarget/100),
				Kind:   "ticket",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "ticket",
			},
		},
	}

	return rules
}

func (*MultiWindowAlgorithm) AlertForLatency(serviceName string, targets []LatencyTarget, annotations map[string]string) []rulefmt.Rule {
	rules := []rulefmt.Rule{
		{
			Alert: "slo:" + serviceName + ".latency.page",
			Expr: multiBurnRateLatency(MultiRateLatencyOpts{
				Metric:  "slo:service_latency",
				Label:   labels.Label{"service", serviceName},
				Buckets: targets,
				Kind:    "page",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "page",
			},
		},
		{
			Alert: "slo:" + serviceName + ".latency.ticket",
			Expr: multiBurnRateLatency(MultiRateLatencyOpts{
				Metric:  "slo:service_latency",
				Label:   labels.Label{"service", serviceName},
				Buckets: targets,
				Kind:    "ticket",
			}),
			Annotations: annotations,
			Labels: map[string]string{
				"severity": "ticket",
			},
		},
	}

	return rules
}

type MultiRateErrorOpts struct {
	Metric string
	Labels labels.Labels
	Value  float64
	Kind   string // page or ticket
}

type MultiRateLatencyOpts struct {
	Metric  string
	Label   labels.Label
	Buckets []LatencyTarget
	Kind    string // page or ticket
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

func multiBurnRate(opts MultiRateErrorOpts) string {
	multiRateWindow := multiRateWindows[opts.Kind]
	conditions := []string{"", ""}

	for index, window := range multiRateWindow {
		condition := fmt.Sprintf(`(%s:ratio_rate_%s%s > (%g * %.3g) and `, opts.Metric, window.LongWindow, opts.Labels.String(), window.Multiplier, opts.Value)
		condition += fmt.Sprintf(`%s:ratio_rate_%s%s > (%g * %.3g))`, opts.Metric, window.ShortWindow, opts.Labels.String(), window.Multiplier, opts.Value)

		conditions[index] = condition
	}

	return strings.Join(conditions, " or ")
}

func multiBurnRateLatency(opts MultiRateLatencyOpts) string {
	multiRateWindow := multiRateWindows[opts.Kind]
	conditions := []string{}

	for _, bucket := range opts.Buckets {
		for _, window := range multiRateWindow {
			value := (1 - ((100 - bucket.Target) * window.Multiplier * 0.01))

			lbs := labels.New(opts.Label, labels.Label{"le", bucket.LE})

			condition := fmt.Sprintf(`(%s:ratio_rate_%s%s < %.3g and `, opts.Metric, window.LongWindow, lbs.String(), value)
			condition += fmt.Sprintf(`%s:ratio_rate_%s%s < %.3g)`, opts.Metric, window.ShortWindow, lbs.String(), value)

			conditions = append(conditions, condition)
		}
	}

	return strings.Join(conditions, " or ")
}

var _ = register(&MultiWindowAlgorithm{}, "multi-window")
