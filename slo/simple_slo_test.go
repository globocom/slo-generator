package slo

import (
	"testing"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	"github.com/stretchr/testify/assert"

	"github.com/globocom/slo-generator/methods"
)

func TestSimpleSLOGenerateAlertRules(t *testing.T) {
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
			BurnRate:    2,
			AlertMethod: "simple",
			AlertWindow: "1h",
			AlertWait:   "5m",
			Expr:        "kk",
		},
		LatencyRecord: ExprBlock{
			AlertMethod: "simple",
			AlertWindow: "30m",
			AlertWait:   "2m",
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

	alertRules := slo.GenerateAlertRules(nil, false)
	assert.Len(t, alertRules, 2)

	assert.Equal(t, ruleNode(rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.errors.page",
		Expr:  "slo:service_errors_total:ratio_rate_1h{service=\"my-team.my-service.payment\"} > 2 * 0.001",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
			"signal":   "error",
		},
		For:         model.Duration(time.Second * 5 * 60),
		Annotations: slo.Annotations,
	}), alertRules[0])

	assert.Equal(t, ruleNode(rulefmt.Rule{
		Alert: "slo:my-team.my-service.payment.latency.page",
		Expr:  "slo:service_latency:ratio_rate_30m{le=\"0.1\", service=\"my-team.my-service.payment\"} < 2 * 0.95 or slo:service_latency:ratio_rate_30m{le=\"0.5\", service=\"my-team.my-service.payment\"} < 2 * 0.99",
		Labels: map[string]string{
			"channel":  "my-channel",
			"severity": "page",
			"signal":   "latency",
		},
		For:         model.Duration(time.Second * 2 * 60),
		Annotations: slo.Annotations,
	}), alertRules[1])
}

func TestSimpleSLOInvalid(t *testing.T) {
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
			AlertMethod: "simple",
			AlertWindow: "22m",
			AlertWait:   "5m",
			Expr:        "kk",
		},
		LatencyRecord: ExprBlock{
			AlertMethod: "simple",
			AlertWindow: "33m",
			AlertWait:   "2m",
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

	assert.PanicsWithValue(t, "Could not generate alert, err: Sample 22m is not a valid sample, valid samples: 5m,30m,1h,2h,6h,1d,3d", func() {
		slo.GenerateAlertRules(nil, false)
	})
}
