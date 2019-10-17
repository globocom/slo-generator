package slo

import (
	"testing"
	"time"

	methods "github.com/globocom/slo-generator/methods"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"
)

func TestSLOGenerateGroupRules(t *testing.T) {
	slo := &SLO{
		Name: "my-team.my-service.payment",
		Objectives: Objectives{
			Availability: 99.9,
			Latency: []methods.LatencyTarget{
				{
					LE:     "0.1",
					Target: 90,
				},
				{
					LE:     "0.5",
					Target: 99,
				},
			},
		},
		TrafficRateRecord: ExprBlock{
			Expr: "sum(rate(http_total[$window]))",
		},
		ErrorRateRecord: ExprBlock{
			AlertMethod: "multi-window",
			Expr:        "sum(rate(http_errors[$window]))/sum(rate(http_total[$window]))",
		},
		LatencyRecord: ExprBlock{
			AlertMethod: "multi-window",
			Expr:        "sum(rate(http_bucket{le=\"$le\"}[$window]))/sum(rate(http_total[$window]))",
		},
		Labels: map[string]string{
			"team": "team-avengers",
		},
		Annotations: map[string]string{
			"message":   "Service A has lower SLI",
			"link":      "http://wiki.ops/1234",
			"dashboard": "http://grafana.globo.com",
		},
	}

	groupRules := slo.GenerateGroupRules()
	assert.Len(t, groupRules, 3)

	assert.Equal(t, groupRules[0], rulefmt.RuleGroup{
		Name:     "slo:my-team.my-service.payment:short",
		Interval: model.Duration(time.Second * 30),
		Rules: []rulefmt.Rule{
			// 5m
			{
				Record: "slo:service_traffic:ratio_rate_5m",
				Expr:   "sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_5m",
				Expr:   "sum(rate(http_errors[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_5m",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_5m",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.5",
				},
			},

			// 30m
			{
				Record: "slo:service_traffic:ratio_rate_30m",
				Expr:   "sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_30m",
				Expr:   "sum(rate(http_errors[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_30m",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_30m",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
					"team":    "team-avengers",
				},
			},

			// 1h
			{
				Record: "slo:service_traffic:ratio_rate_1h",
				Expr:   "sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_1h",
				Expr:   "sum(rate(http_errors[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
					"team":    "team-avengers",
				},
			},
		},
	})

	assert.Equal(t, groupRules[1], rulefmt.RuleGroup{
		Name:     "slo:my-team.my-service.payment:medium",
		Interval: model.Duration(time.Second * 120),
		Rules: []rulefmt.Rule{
			// 2h
			{
				Record: "slo:service_traffic:ratio_rate_2h",
				Expr:   "sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_2h",
				Expr:   "sum(rate(http_errors[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_2h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_2h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
					"team":    "team-avengers",
				},
			},

			// 6h
			{
				Record: "slo:service_traffic:ratio_rate_6h",
				Expr:   "sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_6h",
				Expr:   "sum(rate(http_errors[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_6h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_6h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"team":    "team-avengers",
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},
		},
	})

	assert.Equal(t, groupRules[2], rulefmt.RuleGroup{
		Name:     "slo:my-team.my-service.payment:daily",
		Interval: model.Duration(time.Second * 300),
		Rules: []rulefmt.Rule{
			// 1d
			{
				Record: "slo:service_traffic:ratio_rate_1d",
				Expr:   "sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_1d",
				Expr:   "sum(rate(http_errors[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1d",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1d",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.5",
				},
			},

			// 3d
			{
				Record: "slo:service_traffic:ratio_rate_3d",
				Expr:   "sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_3d",
				Expr:   "sum(rate(http_errors[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_3d",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_3d",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.5",
				},
			},
		},
	})
}

func TestSLOGenerateGroupRulesWithLatencyQuantile(t *testing.T) {
	slo := &SLO{
		Name:        "auto-discover-services",
		HonorLabels: true,
		LatencyQuantileRecord: ExprBlock{
			Expr: "histogram_quantile($quantile, sum by (le) (rate(http_total[$window])))",
		},
	}

	groupRules := slo.GenerateGroupRules()
	assert.Len(t, groupRules, 3)

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:short",
		Interval: model.Duration(time.Second * 30),
		Rules: []rulefmt.Rule{
			// 5m
			{
				Record: "slo:service_latency:p50_5m",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[5m])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_5m",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[5m])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_5m",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[5m])))",
				Labels: map[string]string{},
			},
			// 30m
			{
				Record: "slo:service_latency:p50_30m",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[30m])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_30m",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[30m])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_30m",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[30m])))",
				Labels: map[string]string{},
			},
			// 1h
			{
				Record: "slo:service_latency:p50_1h",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[1h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_1h",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[1h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_1h",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[1h])))",
				Labels: map[string]string{},
			},
		},
	}, groupRules[0])

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:medium",
		Interval: model.Duration(time.Second * 120),
		Rules: []rulefmt.Rule{
			// 2h
			{
				Record: "slo:service_latency:p50_2h",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[2h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_2h",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[2h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_2h",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[2h])))",
				Labels: map[string]string{},
			},
			// 6h
			{
				Record: "slo:service_latency:p50_6h",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[6h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_6h",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[6h])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_6h",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[6h])))",
				Labels: map[string]string{},
			},
		},
	}, groupRules[1])

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:daily",
		Interval: model.Duration(time.Second * 300),
		Rules: []rulefmt.Rule{
			// 1d
			{
				Record: "slo:service_latency:p50_1d",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[1d])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_1d",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[1d])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_1d",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[1d])))",
				Labels: map[string]string{},
			},

			// 3d
			{
				Record: "slo:service_latency:p50_3d",
				Expr:   "histogram_quantile(0.5, sum by (le) (rate(http_total[3d])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p95_3d",
				Expr:   "histogram_quantile(0.95, sum by (le) (rate(http_total[3d])))",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_latency:p99_3d",
				Expr:   "histogram_quantile(0.99, sum by (le) (rate(http_total[3d])))",
				Labels: map[string]string{},
			},
		},
	}, groupRules[2])
}

func TestSLOGenerateGroupRulesWithAutoDiscovery(t *testing.T) {
	slo := &SLO{
		Name:        "auto-discover-services",
		HonorLabels: true,
		TrafficRateRecord: ExprBlock{
			Expr: "sum(rate(http_total[$window])) by (service)",
		},
		ErrorRateRecord: ExprBlock{
			AlertMethod: "multi-window",
			Expr:        "sum(rate(http_errors[$window])) by (service)/sum(rate(http_total[$window])) by (service)",
		},
	}

	groupRules := slo.GenerateGroupRules()
	assert.Len(t, groupRules, 3)

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:short",
		Interval: model.Duration(time.Second * 30),
		Rules: []rulefmt.Rule{
			// 5m
			{
				Record: "slo:service_traffic:ratio_rate_5m",
				Expr:   "sum(rate(http_total[5m])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_5m",
				Expr:   "sum(rate(http_errors[5m])) by (service)/sum(rate(http_total[5m])) by (service)",
				Labels: map[string]string{},
			},
			// 30m
			{
				Record: "slo:service_traffic:ratio_rate_30m",
				Expr:   "sum(rate(http_total[30m])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_30m",
				Expr:   "sum(rate(http_errors[30m])) by (service)/sum(rate(http_total[30m])) by (service)",
				Labels: map[string]string{},
			},
			// 1h
			{
				Record: "slo:service_traffic:ratio_rate_1h",
				Expr:   "sum(rate(http_total[1h])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_1h",
				Expr:   "sum(rate(http_errors[1h])) by (service)/sum(rate(http_total[1h])) by (service)",
				Labels: map[string]string{},
			},
		},
	}, groupRules[0])

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:medium",
		Interval: model.Duration(time.Second * 120),
		Rules: []rulefmt.Rule{
			// 2h
			{
				Record: "slo:service_traffic:ratio_rate_2h",
				Expr:   "sum(rate(http_total[2h])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_2h",
				Expr:   "sum(rate(http_errors[2h])) by (service)/sum(rate(http_total[2h])) by (service)",
				Labels: map[string]string{},
			},

			// 6h
			{
				Record: "slo:service_traffic:ratio_rate_6h",
				Expr:   "sum(rate(http_total[6h])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_6h",
				Expr:   "sum(rate(http_errors[6h])) by (service)/sum(rate(http_total[6h])) by (service)",
				Labels: map[string]string{},
			},
		},
	}, groupRules[1])

	assert.Equal(t, rulefmt.RuleGroup{
		Name:     "slo:auto-discover-services:daily",
		Interval: model.Duration(time.Second * 300),
		Rules: []rulefmt.Rule{
			// 1d
			{
				Record: "slo:service_traffic:ratio_rate_1d",
				Expr:   "sum(rate(http_total[1d])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_1d",
				Expr:   "sum(rate(http_errors[1d])) by (service)/sum(rate(http_total[1d])) by (service)",
				Labels: map[string]string{},
			},

			// 3d
			{
				Record: "slo:service_traffic:ratio_rate_3d",
				Expr:   "sum(rate(http_total[3d])) by (service)",
				Labels: map[string]string{},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_3d",
				Expr:   "sum(rate(http_errors[3d])) by (service)/sum(rate(http_total[3d])) by (service)",
				Labels: map[string]string{},
			},
		},
	}, groupRules[2])
}

func TestSLOGenerateAlertRules(t *testing.T) {
	slo := &SLO{
		Name: "my-team.my-service.payment",
		Objectives: Objectives{
			Availability: 99.9,
			Latency: []methods.LatencyTarget{
				{
					LE:     "0.1",
					Target: 95,
				},
				{
					LE:     "0.5",
					Target: 99,
				},
			},
		},
		ErrorRateRecord: ExprBlock{
			AlertMethod: "multi-window",
			Expr:        "kk",
		},
		LatencyRecord: ExprBlock{
			AlertMethod: "multi-window",
			Expr:        "kk",
		},
		Labels: map[string]string{
			"channel": "my-channel",
		},
		Annotations: map[string]string{
			"message":   "Service A has lower SLI",
			"link":      "http://wiki.ops/1234",
			"dashboard": "http://grafana.globo.com",
		},
	}

	alertRules := slo.GenerateAlertRules()
	assert.Len(t, alertRules, 4)

	assert.Equal(t, alertRules[0], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.page",
		Expr:  "(slo:service_errors_total:ratio_rate_1h{service=\"my-team.my-service.payment\"} > (14.4 * 0.001) and slo:service_errors_total:ratio_rate_5m{service=\"my-team.my-service.payment\"} > (14.4 * 0.001)) or (slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (6 * 0.001) and slo:service_errors_total:ratio_rate_30m{service=\"my-team.my-service.payment\"} > (6 * 0.001))",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[1], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.ticket",
		Expr:  "(slo:service_errors_total:ratio_rate_1d{service=\"my-team.my-service.payment\"} > (3 * 0.001) and slo:service_errors_total:ratio_rate_2h{service=\"my-team.my-service.payment\"} > (3 * 0.001)) or (slo:service_errors_total:ratio_rate_3d{service=\"my-team.my-service.payment\"} > (1 * 0.001) and slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (1 * 0.001))",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "ticket",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[2], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.latency.page",
		Expr: ("(" +
			"slo:service_latency:ratio_rate_1h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.28" +
			" and " +
			"slo:service_latency:ratio_rate_5m{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.28" +
			") or (" +
			"slo:service_latency:ratio_rate_6h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.7" +
			" and " +
			"slo:service_latency:ratio_rate_30m{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.7" +
			") or (" +
			"slo:service_latency:ratio_rate_1h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.856" +
			" and " +
			"slo:service_latency:ratio_rate_5m{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.856" +
			") or (" +
			"slo:service_latency:ratio_rate_6h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.94" +
			" and " +
			"slo:service_latency:ratio_rate_30m{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.94" +
			")"),

		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[3], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.latency.ticket",
		Expr: ("(" +
			"slo:service_latency:ratio_rate_1d{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.85" +
			" and " +
			"slo:service_latency:ratio_rate_2h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.85" +
			") or (" +
			"slo:service_latency:ratio_rate_3d{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.95" +
			" and " +
			"slo:service_latency:ratio_rate_6h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.95" +
			") or (" +
			"slo:service_latency:ratio_rate_1d{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.97" +
			" and " +
			"slo:service_latency:ratio_rate_2h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.97" +
			") or (" +
			"slo:service_latency:ratio_rate_3d{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.99" +
			" and " +
			"slo:service_latency:ratio_rate_6h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.99" +
			")"),

		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "ticket",
		},
		Annotations: slo.Annotations,
	})
}

func TestSLOGenerateAlertRulesWithoutExpressions(t *testing.T) {
	slo := &SLO{
		Name: "my-team.my-service.payment",
		Objectives: Objectives{
			Availability: 99.9,
			Latency: []methods.LatencyTarget{
				{
					LE:     "0.1",
					Target: 95,
				},
				{
					LE:     "0.5",
					Target: 99,
				},
			},
		},
		ErrorRateRecord: ExprBlock{
			AlertMethod: "multi-window",
		},
		LatencyRecord: ExprBlock{
			AlertMethod: "multi-window",
		},
		Labels: map[string]string{
			"channel": "my-channel",
		},
		Annotations: map[string]string{
			"message":   "Service A has lower SLI",
			"link":      "http://wiki.ops/1234",
			"dashboard": "http://grafana.globo.com",
		},
	}

	alertRules := slo.GenerateAlertRules()
	assert.Len(t, alertRules, 4)

	assert.Equal(t, alertRules[0], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.page",
		Expr:  "(slo:service_errors_total:ratio_rate_1h{service=\"my-team.my-service.payment\"} > (14.4 * 0.001) and slo:service_errors_total:ratio_rate_5m{service=\"my-team.my-service.payment\"} > (14.4 * 0.001)) or (slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (6 * 0.001) and slo:service_errors_total:ratio_rate_30m{service=\"my-team.my-service.payment\"} > (6 * 0.001))",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[1], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.ticket",
		Expr:  "(slo:service_errors_total:ratio_rate_1d{service=\"my-team.my-service.payment\"} > (3 * 0.001) and slo:service_errors_total:ratio_rate_2h{service=\"my-team.my-service.payment\"} > (3 * 0.001)) or (slo:service_errors_total:ratio_rate_3d{service=\"my-team.my-service.payment\"} > (1 * 0.001) and slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (1 * 0.001))",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "ticket",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[2], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.latency.page",
		Expr: ("(" +
			"slo:service_latency:ratio_rate_1h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.28" +
			" and " +
			"slo:service_latency:ratio_rate_5m{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.28" +
			") or (" +
			"slo:service_latency:ratio_rate_6h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.7" +
			" and " +
			"slo:service_latency:ratio_rate_30m{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.7" +
			") or (" +
			"slo:service_latency:ratio_rate_1h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.856" +
			" and " +
			"slo:service_latency:ratio_rate_5m{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.856" +
			") or (" +
			"slo:service_latency:ratio_rate_6h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.94" +
			" and " +
			"slo:service_latency:ratio_rate_30m{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.94" +
			")"),

		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[3], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.latency.ticket",
		Expr: ("(" +
			"slo:service_latency:ratio_rate_1d{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.85" +
			" and " +
			"slo:service_latency:ratio_rate_2h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.85" +
			") or (" +
			"slo:service_latency:ratio_rate_3d{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.95" +
			" and " +
			"slo:service_latency:ratio_rate_6h{le=\"0.1\", service=\"my-team.my-service.payment\"} < 0.95" +
			") or (" +
			"slo:service_latency:ratio_rate_1d{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.97" +
			" and " +
			"slo:service_latency:ratio_rate_2h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.97" +
			") or (" +
			"slo:service_latency:ratio_rate_3d{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.99" +
			" and " +
			"slo:service_latency:ratio_rate_6h{le=\"0.5\", service=\"my-team.my-service.payment\"} < 0.99" +
			")"),

		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "ticket",
		},
		Annotations: slo.Annotations,
	})
}
