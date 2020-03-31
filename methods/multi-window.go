package methods

import (
	"fmt"
	"strings"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type MultiWindowAlgorithm struct{}

func (*MultiWindowAlgorithm) AlertForError(opts *AlertErrorOptions) ([]rulefmt.Rule, error) {
	rules := []rulefmt.Rule{
		{
			Alert: "slo:" + opts.ServiceName + ".errors.page",
			Expr: multiBurnRate(MultiRateErrorOpts{
				Metric: "slo:service_errors_total",
				Labels: labels.New(labels.Label{"service", opts.ServiceName}),
				Value:  (1 - opts.AvailabilityTarget/100),
				Kind:   "page",
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "page",
			},
		},
		{
			Alert: "slo:" + opts.ServiceName + ".errors.ticket",
			Expr: multiBurnRate(MultiRateErrorOpts{
				Metric: "slo:service_errors_total",
				Labels: labels.New(labels.Label{"service", opts.ServiceName}),
				Value:  (1 - opts.AvailabilityTarget/100),
				Kind:   "ticket",
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "ticket",
			},
		},
	}

	return rules, nil
}

func (*MultiWindowAlgorithm) AlertForLatency(opts *AlertLatencyOptions) ([]rulefmt.Rule, error) {
	rules := []rulefmt.Rule{
		{
			Alert: "slo:" + opts.ServiceName + ".latency.page",
			Expr: multiBurnRateLatency(MultiRateLatencyOpts{
				Metric:  "slo:service_latency",
				Label:   labels.Label{"service", opts.ServiceName},
				Buckets: opts.Targets,
				Kind:    "page",
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "page",
			},
		},
		{
			Alert: "slo:" + opts.ServiceName + ".latency.ticket",
			Expr: multiBurnRateLatency(MultiRateLatencyOpts{
				Metric:  "slo:service_latency",
				Label:   labels.Label{"service", opts.ServiceName},
				Buckets: opts.Targets,
				Kind:    "ticket",
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": "ticket",
			},
		},
	}

	return rules, nil
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
