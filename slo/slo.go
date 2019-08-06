package slo

import (
	"log"
	"strings"

	algorithms "github.com/globocom/slo-generator/algorithms"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type SLOSpec struct {
	SLOS []SLO
}

type ExprBlock struct {
	AlertAlgorithm string `yaml:"alertAlgorithm"`
	Expr           string `yaml:"expr"`
}

func (block *ExprBlock) ComputeExpr(window, le string) string {
	replacer := strings.NewReplacer("$window", window, "$le", le)
	return replacer.Replace(block.Expr)
}

type SLO struct {
	Name                         string            `yaml:"name"`
	AvailabilityObjectivePercent float64           `yaml:"availabilityObjectivePercent"`
	LatencyObjectiveBuckets      []LatencyBucket   `yaml:"latencyObjectiveBuckets"`
	ErrorRateRecord              ExprBlock         `yaml:"errorRateRecord"`
	LatencyRecord                ExprBlock         `yaml:"latencyRecord"`
	Annotations                  map[string]string `yaml:"annotations"`
}

type LatencyBucket struct {
	LE     string  `yaml:"le"`
	Target float64 `yaml:"target"`
}

func (slo SLO) GenerateAlertRules() []rulefmt.Rule {
	alertRules := []rulefmt.Rule{}

	errorAlgorithm := algorithms.Get(slo.ErrorRateRecord.AlertAlgorithm)
	if errorAlgorithm != nil {
		errorRules := errorAlgorithm.AlertForError(slo.Name, slo.AvailabilityObjectivePercent, slo.Annotations)
		alertRules = append(alertRules, errorRules...)
	}

	latencyAlgorithm := algorithms.Get(slo.LatencyRecord.AlertAlgorithm)
	if latencyAlgorithm != nil {

	}
	// if slo.Algorithm == "multiwindow" {
	//	// alerting page
	//	sloPageRecord := rulefmt.Rule{
	//		Alert: "slo:" + slo.Name + ".errors.page",
	//		Expr: algorithms.MultiBurnRateForPage(
	//			"slo:service_errors_total",
	//			labels.New(labels.Label{"service", slo.Name}),
	//			">", (1 - slo.AvailabilityObjectivePercent/100),
	//		),
	//		Annotations: slo.Annotations,
	//		Labels: map[string]string{
	//			"severity": "page",
	//		},
	//	}

	//	alertRules = append(alertRules, sloPageRecord)

	//	// alerting ticket
	//	sloTicketRecord := rulefmt.Rule{
	//		Alert: "slo:" + slo.Name + ".errors.ticket",
	//		Expr: algorithms.MultiBurnRateForTicket(
	//			"slo:service_errors_total",
	//			labels.New(labels.Label{"service", slo.Name}),
	//			">", (1 - slo.AvailabilityObjectivePercent/100),
	//		),
	//		Annotations: slo.Annotations,
	//		Labels: map[string]string{
	//			"severity": "ticket",
	//		},
	//	}

	//	alertRules = append(alertRules, sloTicketRecord)

	// }

	return alertRules
}

func (slo SLO) GenerateGroupRules() []rulefmt.RuleGroup {
	rules := []rulefmt.RuleGroup{}

	for _, sample := range defaultSamples {
		interval, err := model.ParseDuration(sample.Interval)
		if err != nil {
			log.Fatal(err)
		}
		ruleGroup := rulefmt.RuleGroup{
			Name:     "slo:" + slo.Name + ":" + sample.Name,
			Interval: interval,
			Rules:    []rulefmt.Rule{},
		}

		for _, bucket := range sample.Buckets {
			errorRateRecord := rulefmt.Rule{
				Record: "slo:service_errors_total:ratio_rate_" + bucket,
				Expr:   slo.ErrorRateRecord.ComputeExpr(bucket, ""),
				Labels: map[string]string{
					"service": slo.Name,
				},
			}

			ruleGroup.Rules = append(ruleGroup.Rules, errorRateRecord)

			for _, latencyBucket := range slo.LatencyObjectiveBuckets {
				latencyRateRecord := rulefmt.Rule{
					Record: "slo:service_latency:ratio_rate_" + bucket,
					Expr:   slo.LatencyRecord.ComputeExpr(bucket, latencyBucket.LE),
					Labels: map[string]string{
						"service": slo.Name,
						"le":      latencyBucket.LE,
					},
				}

				ruleGroup.Rules = append(ruleGroup.Rules, latencyRateRecord)
			}

		}

		rules = append(rules, ruleGroup)
	}

	return rules
}
