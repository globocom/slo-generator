package kubernetes

import (
	"testing"

	"github.com/globocom/slo-generator/methods"
	"github.com/globocom/slo-generator/slo"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestGenerateManifests(t *testing.T) {
	manifests := GenerateManifests(Opts{
		SLO: slo.SLO{
			Name: "my-team.my-service.payment",
			Objectives: slo.Objectives{
				Availability: 99.9,
				Latency: []methods.LatencyTarget{
					{
						LE:     "0.5",
						Target: 99,
					},
				},
			},
			TrafficRateRecord: slo.ExprBlock{
				Expr: "sum(rate(http_total[$window]))",
			},
			ErrorRateRecord: slo.ExprBlock{
				AlertMethod: "multi-window",
				Expr:        "sum(rate(http_errors[$window]))/sum(rate(http_total[$window]))",
			},
			LatencyRecord: slo.ExprBlock{
				AlertMethod: "multi-window",
				Expr:        "sum(rate(http_bucket{le=\"$le\"}[$window]))/sum(rate(http_total[$window]))",
			},
			Labels: map[string]string{
				"team": "team-avengers",
			},
			Annotations: map[string]string{
				"message": "Service A has lower SLI",
			},
		},
	})

	assert.Len(t, manifests, 2)
	assert.Equal(t, v1.ObjectMeta{
		Name: "slis-my-team.my-service.payment",
	}, manifests[0].ObjectMeta)

	assert.Equal(t, monitoringv1.RuleGroup{
		Name:     "slo:my-team.my-service.payment:short",
		Interval: "30s",
		Rules: []monitoringv1.Rule{
			// 5m
			{
				Record: "slo:service_traffic:ratio_rate_5m",
				Expr:   intstr.FromString("sum(rate(http_total[5m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_5m",
				Expr:   intstr.FromString("sum(rate(http_errors[5m]))/sum(rate(http_total[5m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_5m",
				Expr:   intstr.FromString("sum(rate(http_bucket{le=\"0.5\"}[5m]))/sum(rate(http_total[5m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
					"le":      "0.5",
				},
			},

			// 30m
			{
				Record: "slo:service_traffic:ratio_rate_30m",
				Expr:   intstr.FromString("sum(rate(http_total[30m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_30m",
				Expr:   intstr.FromString("sum(rate(http_errors[30m]))/sum(rate(http_total[30m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_30m",
				Expr:   intstr.FromString("sum(rate(http_bucket{le=\"0.5\"}[30m]))/sum(rate(http_total[30m]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
					"team":    "team-avengers",
				},
			},

			// 1h
			{
				Record: "slo:service_traffic:ratio_rate_1h",
				Expr:   intstr.FromString("sum(rate(http_total[1h]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_errors_total:ratio_rate_1h",
				Expr:   intstr.FromString("sum(rate(http_errors[1h]))/sum(rate(http_total[1h]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"team":    "team-avengers",
				},
			},
			{
				Record: "slo:service_latency:ratio_rate_1h",
				Expr:   intstr.FromString("sum(rate(http_bucket{le=\"0.5\"}[1h]))/sum(rate(http_total[1h]))"),
				Labels: map[string]string{
					"service": "my-team.my-service.payment",
					"le":      "0.5",
					"team":    "team-avengers",
				},
			},
		},
	}, manifests[0].Spec.Groups[0])

	assert.Equal(t, v1.ObjectMeta{
		Name: "slos-alerts-my-team.my-service.payment",
	}, manifests[1].ObjectMeta)
}
