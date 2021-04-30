package methods

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type MultiWindowAlgorithm struct{}

func (*MultiWindowAlgorithm) AlertForError(opts *AlertErrorOptions) ([]rulefmt.Rule, error) {
	ratesMap := genMultiRateWindows(opts.SLOWindow, opts.ShortWindow, opts.Windows)
	rules := []rulefmt.Rule{}

	for _, severity := range Severities {
		if _, ok := ratesMap[severity]; !ok {
			continue
		}
		rules = append(rules, rulefmt.Rule{
			Alert: "slo:" + opts.ServiceName + ".errors." + string(severity),
			Expr: multiBurnRate(MultiRateErrorOpts{
				Rates:  ratesMap[severity],
				Metric: "slo:service_errors_total",
				Labels: labels.New(labels.Label{Name: "service", Value: opts.ServiceName}),
				Value:  1 - opts.AvailabilityTarget/100,
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": string(severity),
				"signal":   "error",
			},
		})
	}
	return rules, nil
}

func (*MultiWindowAlgorithm) AlertForLatency(opts *AlertLatencyOptions) ([]rulefmt.Rule, error) {
	ratesMap := genMultiRateWindows(opts.SLOWindow, opts.ShortWindow, opts.Windows)
	rules := []rulefmt.Rule{}

	for _, severity := range Severities {
		if _, ok := ratesMap[severity]; !ok {
			continue
		}
		rules = append(rules, rulefmt.Rule{
			Alert: "slo:" + opts.ServiceName + ".latency." + string(severity),
			Expr: multiBurnRateLatency(MultiRateLatencyOpts{
				Rates:   ratesMap[severity],
				Metric:  "slo:service_latency",
				Label:   labels.Label{Name: "service", Value: opts.ServiceName},
				Buckets: opts.Targets,
			}),
			Annotations: map[string]string{},
			Labels: map[string]string{
				"severity": string(severity),
				"signal":   "latency",
			},
		})
	}

	return rules, nil
}

type MultiRateErrorOpts struct {
	Rates  []MultiRateWindow
	Metric string
	Labels labels.Labels
	Value  float64
}

type MultiRateLatencyOpts struct {
	Rates   []MultiRateWindow
	Metric  string
	Label   labels.Label
	Buckets []LatencyTarget
}

type MultiRateWindow struct {
	Multiplier  float64
	LongWindow  string
	ShortWindow string
}

var multiRateWindows = map[NotificationSeverity][]MultiRateWindow{
	NotificationPageSeverity: {
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
	NotificationTicketSeverity: {
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

func genMultiRateWindows(SLOWindow time.Duration, shortWindow bool, windows []Window) map[NotificationSeverity][]MultiRateWindow {
	if len(windows) == 0 {
		// Use Default multiRateWindows from SRE Book
		return multiRateWindows
	}

	mrate := map[NotificationSeverity][]MultiRateWindow{}
	wHours := float64(SLOWindow / time.Hour)

	for _, w := range windows {
		t := float64(time.Duration(w.Duration) / time.Hour)

		burnRate := (w.Consumption / 100) / (t / wHours)
		m := MultiRateWindow{
			Multiplier: burnRate,
			LongWindow: w.Duration.String(),
		}

		if shortWindow {
			// Short window is defined as 1/12 of the long window for now
			short := time.Duration(w.Duration) / 12
			m.ShortWindow = model.Duration(short).String()
		}
		mrate[w.Notification] = append(mrate[w.Notification], m)
	}

	return mrate
}

func multiBurnRate(opts MultiRateErrorOpts) string {
	multiRateWindow := opts.Rates
	conditions := []string{}

	for _, window := range multiRateWindow {
		condition := fmt.Sprintf(`%s:ratio_rate_%s%s > (%g * %.3g)`, opts.Metric, window.LongWindow, opts.Labels.String(), window.Multiplier, opts.Value)
		if window.ShortWindow != "" {
			condition = fmt.Sprintf(`(%s and %s:ratio_rate_%s%s > (%g * %.3g))`, condition, opts.Metric, window.ShortWindow, opts.Labels.String(), window.Multiplier, opts.Value)
		}

		conditions = append(conditions, condition)
	}

	if len(conditions) == 1 {
		return conditions[0]
	}

	return strings.Join(conditions, " or ")
}

func multiBurnRateLatency(opts MultiRateLatencyOpts) string {
	multiRateWindow := opts.Rates
	conditions := []string{}

	for _, bucket := range opts.Buckets {
		for _, window := range multiRateWindow {
			value := (1 - ((100 - bucket.Target) / 100 * window.Multiplier))
			lbs := labels.New(opts.Label, labels.Label{Name: "le", Value: bucket.LE})

			condition := fmt.Sprintf(`%s:ratio_rate_%s%s < %.3g`, opts.Metric, window.LongWindow, lbs.String(), value)
			if window.ShortWindow != "" {
				condition = fmt.Sprintf(`(%s and %s:ratio_rate_%s%s < %.3g)`, condition, opts.Metric, window.ShortWindow, lbs.String(), value)
			}

			conditions = append(conditions, condition)
		}
	}

	if len(conditions) == 1 {
		return conditions[0]
	}

	return strings.Join(conditions, " or ")
}

var _ = register(&MultiWindowAlgorithm{}, "multi-window")
