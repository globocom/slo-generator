package slo

import (
	"log"
	"strings"

	methods "github.com/globocom/slo-generator/methods"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/rulefmt"
)

type SLOSpec struct {
	SLOS []SLO
}

type ExprBlock struct {
	AlertMethod string `yaml:"alertMethod"`
	Expr        string `yaml:"expr"`
}

func (block *ExprBlock) ComputeExpr(window, le string) string {
	replacer := strings.NewReplacer("$window", window, "$le", le)
	return replacer.Replace(block.Expr)
}

type SLO struct {
	Name       string `yaml:"name"`
	Objectives Objectives

	HonorLabels bool `yaml:"honorLabels"`

	TrafficRateRecord ExprBlock         `yaml:"trafficRateRecord"`
	ErrorRateRecord   ExprBlock         `yaml:"errorRateRecord"`
	LatencyRecord     ExprBlock         `yaml:"latencyRecord"`
	Labels            map[string]string `yaml:"labels"`
	Annotations       map[string]string `yaml:"annotations"`
}

type Objectives struct {
	Availability float64                 `yaml:"availability"`
	Latency      []methods.LatencyTarget `yaml:"latency"`
}

func (slo SLO) GenerateAlertRules() []rulefmt.Rule {
	alertRules := []rulefmt.Rule{}

	errorMethod := methods.Get(slo.ErrorRateRecord.AlertMethod)
	if errorMethod != nil {
		errorRules := errorMethod.AlertForError(slo.Name, slo.Objectives.Availability)
		alertRules = append(alertRules, errorRules...)
	}

	latencyMethod := methods.Get(slo.LatencyRecord.AlertMethod)
	if latencyMethod != nil {
		latencyRules := latencyMethod.AlertForLatency(slo.Name, slo.Objectives.Latency)
		alertRules = append(alertRules, latencyRules...)
	}

	for _, rule := range alertRules {
		slo.fillMetadata(&rule)
	}

	return alertRules
}

func (slo *SLO) fillMetadata(rule *rulefmt.Rule) {
	for label, value := range slo.Labels {
		rule.Labels[label] = value
	}

	for label, value := range slo.Annotations {
		rule.Annotations[label] = value
	}
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
			ruleGroup.Rules = append(ruleGroup.Rules, slo.generateRules(bucket)...)
		}

		rules = append(rules, ruleGroup)
	}

	return rules
}

func (slo SLO) generateRules(bucket string) []rulefmt.Rule {
	rules := []rulefmt.Rule{}
	if slo.TrafficRateRecord.Expr != "" {
		trafficRateRecord := rulefmt.Rule{
			Record: "slo:service_traffic:ratio_rate_" + bucket,
			Expr:   slo.TrafficRateRecord.ComputeExpr(bucket, ""),
			Labels: map[string]string{},
		}

		if !slo.HonorLabels {
			trafficRateRecord.Labels["service"] = slo.Name
		}

		rules = append(rules, trafficRateRecord)
	}

	errorRateRecord := rulefmt.Rule{
		Record: "slo:service_errors_total:ratio_rate_" + bucket,
		Expr:   slo.ErrorRateRecord.ComputeExpr(bucket, ""),
		Labels: map[string]string{},
	}

	if !slo.HonorLabels {
		errorRateRecord.Labels["service"] = slo.Name
	}

	rules = append(rules, errorRateRecord)

	for _, latencyBucket := range slo.Objectives.Latency {
		latencyRateRecord := rulefmt.Rule{
			Record: "slo:service_latency:ratio_rate_" + bucket,
			Expr:   slo.LatencyRecord.ComputeExpr(bucket, latencyBucket.LE),
			Labels: map[string]string{
				"service": slo.Name,
				"le":      latencyBucket.LE,
			},
		}

		rules = append(rules, latencyRateRecord)
	}

	return rules
}
