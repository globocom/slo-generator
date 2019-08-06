package slo

import (
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"
)

func TestSLOGenerateGroupRules(t *testing.T) {
	slo := &SLO{
		Name: "my-team.my-service.payment",
		AvailabilityObjectivePercent: 99.9,
		LatencyObjectiveBuckets: []LatencyBucket{
			{
				LE:     "0.1",
				Target: 90,
			},
			{
				LE:     "0.5",
				Target: 99,
			},
		},
		ErrorRateRecord: ExprBlock{
			AlertAlgorithm: "multi-window",
			Expr:           "sum(rate(http_errors[$window]))/sum(rate(http_total[$window]))",
		},
		LatencyRecord: ExprBlock{
			AlertAlgorithm: "multi-window",
			Expr:           "sum(rate(http_bucket{le=\"$le\"}[$window]))/sum(rate(http_total[$window]))",
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
				Record: "slo:service_errors_total:ratio_rate_5m",
				Expr:   "sum(rate(http_errors[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_5m",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_5m",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[5m]))/sum(rate(http_total[5m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},

			// 30m
			{
				Record: "slo:service_errors_total:ratio_rate_30m",
				Expr:   "sum(rate(http_errors[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_30m",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_30m",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[30m]))/sum(rate(http_total[30m]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},

			// 1h
			{
				Record: "slo:service_errors_total:ratio_rate_1h",
				Expr:   "sum(rate(http_errors[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[1h]))/sum(rate(http_total[1h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
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
				Record: "slo:service_errors_total:ratio_rate_2h",
				Expr:   "sum(rate(http_errors[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_2h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_2h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[2h]))/sum(rate(http_total[2h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},

			// 6h
			{
				Record: "slo:service_errors_total:ratio_rate_6h",
				Expr:   "sum(rate(http_errors[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_6h",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_6h",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[6h]))/sum(rate(http_total[6h]))",
				Labels: map[string]string{
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
				Record: "slo:service_errors_total:ratio_rate_1d",
				Expr:   "sum(rate(http_errors[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1d",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1d",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[1d]))/sum(rate(http_total[1d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},

			// 3d
			{
				Record: "slo:service_errors_total:ratio_rate_3d",
				Expr:   "sum(rate(http_errors[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_3d",
				Expr:   "sum(rate(http_bucket{le=\"0.1\"}[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.1",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_3d",
				Expr:   "sum(rate(http_bucket{le=\"0.5\"}[3d]))/sum(rate(http_total[3d]))",
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
				},
			},
		},
	})
}

func TestSLOGenerateAlertRules(t *testing.T) {
	slo := &SLO{
		Name: "my-team.my-service.payment",
		AvailabilityObjectivePercent: 99.9,
		LatencyObjectiveBuckets: []LatencyBucket{
			{
				LE:     "0.1",
				Target: 90,
			},
			{
				LE:     "0.5",
				Target: 99,
			},
		},
		ErrorRateRecord: ExprBlock{
			AlertAlgorithm: "multi-window",
			Expr:           "kk",
		},
		LatencyRecord: ExprBlock{
			AlertAlgorithm: "multi-window",
			Expr:           "kk",
		},
		Annotations: map[string]string{
			"message":   "Service A has lower SLI",
			"link":      "http://wiki.ops/1234",
			"dashboard": "http://grafana.globo.com",
		},
	}

	alertRules := slo.GenerateAlertRules()
	assert.Len(t, alertRules, 2)

	assert.Equal(t, alertRules[0], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.page",
		Expr:  "(slo:service_errors_total:ratio_rate_1h{service=\"my-team.my-service.payment\"} > (14.4 * 0.001) and slo:service_errors_total:ratio_rate_5m{service=\"my-team.my-service.payment\"} > (14.4 * 0.001)) or (slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (6 * 0.001) and slo:service_errors_total:ratio_rate_30m{service=\"my-team.my-service.payment\"} > (6 * 0.001))",
		Labels: map[string]string{
			"severity": "page",
		},
		Annotations: slo.Annotations,
	})

	assert.Equal(t, alertRules[1], rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.ticket",
		Expr:  "(slo:service_errors_total:ratio_rate_1d{service=\"my-team.my-service.payment\"} > (3 * 0.001) and slo:service_errors_total:ratio_rate_2h{service=\"my-team.my-service.payment\"} > (3 * 0.001)) or (slo:service_errors_total:ratio_rate_3d{service=\"my-team.my-service.payment\"} > (1 * 0.001) and slo:service_errors_total:ratio_rate_6h{service=\"my-team.my-service.payment\"} > (1 * 0.001))",
		Labels: map[string]string{
			"severity": "ticket",
		},
		Annotations: slo.Annotations,
	})
}
