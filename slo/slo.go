package slo

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/globocom/slo-generator/methods"
	"github.com/globocom/slo-generator/samples"
	"github.com/prometheus/common/model"

	"github.com/prometheus/prometheus/pkg/rulefmt"
)

var quantiles = []struct {
	name     string
	quantile float64
}{
	{
		name:     "p50",
		quantile: 0.5,
	},
	{
		name:     "p95",
		quantile: 0.95,
	},
	{
		name:     "p99",
		quantile: 0.99,
	},
}

type SLOSpec struct {
	SLOS    []SLO   `yaml:"slos"`
	Classes Classes `yaml:"classes"`
}

type ExprBlock struct {
	AlertMethod string           `yaml:"alertMethod"`
	AlertWindow string           `yaml:"alertWindow"`
	BurnRate    float64          `yaml:"burnRate"`
	AlertWait   string           `yaml:"alertWait"`
	Windows     []methods.Window `yaml:"windows"`
	ShortWindow *bool            `yaml:"shortWindow"`
	Buckets     []string         `yaml:"buckets"` // used to define buckets of histogram when using latency expression
	Expr        string           `yaml:"expr"`
}

func (block *ExprBlock) GetShortWindow() bool {
	defaultShortWindow := true

	if block.ShortWindow == nil {
		return defaultShortWindow
	}

	return *block.ShortWindow
}

func (block *ExprBlock) ComputeExpr(window, le string) string {
	replacer := strings.NewReplacer("$window", window, "$le", le)
	return replacer.Replace(block.Expr)
}

func (block *ExprBlock) ComputeQuantile(window string, quantile float64) string {
	replacer := strings.NewReplacer("$window", window, "$quantile", fmt.Sprintf("%g", quantile))
	return replacer.Replace(block.Expr)
}

type SLO struct {
	Name       string `yaml:"name"`
	Class      string `yaml:"class"`
	Objectives Objectives

	HonorLabels bool `yaml:"honorLabels"`

	TrafficRateRecord     ExprBlock         `yaml:"trafficRateRecord"`
	ErrorRateRecord       ExprBlock         `yaml:"errorRateRecord"`
	LatencyRecord         ExprBlock         `yaml:"latencyRecord"`
	LatencyQuantileRecord ExprBlock         `yaml:"latencyQuantileRecord"`
	Labels                map[string]string `yaml:"labels"`
	Annotations           map[string]string `yaml:"annotations"`
}

type Objectives struct {
	Availability float64                 `yaml:"availability"`
	Latency      []methods.LatencyTarget `yaml:"latency"`
	Window       model.Duration          `yaml:"window"`
}

// LatencyBuckets returns all boundaries of latencies
// is the same boundaries of a prometheus histogram (aka: le) used to calculate latency SLOs
func (o *Objectives) LatencyBuckets() []string {
	var latencyBuckets []string

	for _, latencyBucket := range o.Latency {
		latencyBuckets = append(latencyBuckets, latencyBucket.LE)
	}

	return latencyBuckets
}

func (slo *SLO) GenerateAlertRules(sloClass *Class, disableTicket bool) []rulefmt.RuleNode {
	objectives := slo.Objectives
	if sloClass != nil {
		objectives = sloClass.Objectives
	}

	var alertRules []rulefmt.RuleNode

	if slo.ErrorRateRecord.AlertMethod != "" {
		errorMethod := methods.Get(slo.ErrorRateRecord.AlertMethod)
		if errorMethod == nil {
			log.Panicf("alertMethod %s is not valid", slo.ErrorRateRecord.AlertMethod)
		}

		errorRules, err := errorMethod.AlertForError(&methods.AlertErrorOptions{
			ServiceName:        slo.Name,
			AvailabilityTarget: objectives.Availability,
			SLOWindow:          time.Duration(objectives.Window),
			ShortWindow:        slo.ErrorRateRecord.GetShortWindow(),
			Windows:            slo.ErrorRateRecord.Windows,
			AlertWindow:        slo.ErrorRateRecord.AlertWindow,
			AlertWait:          slo.ErrorRateRecord.AlertWait,
			BurnRate:           slo.ErrorRateRecord.BurnRate,
		})
		if err != nil {
			log.Panicf("Could not generate alert, err: %s", err.Error())
		}
		alertRules = append(alertRules, ruleNodes(errorRules)...)
	}

	if slo.LatencyRecord.AlertMethod != "" {
		latencyMethod := methods.Get(slo.LatencyRecord.AlertMethod)
		if latencyMethod == nil {
			log.Panicf("alertMethod %s is not valid", slo.LatencyRecord.AlertMethod)
		}

		if objectives.Latency != nil {
			latencyRules, err := latencyMethod.AlertForLatency(&methods.AlertLatencyOptions{
				ServiceName: slo.Name,
				Targets:     objectives.Latency,
				SLOWindow:   time.Duration(objectives.Window),
				ShortWindow: slo.LatencyRecord.GetShortWindow(),
				Windows:     slo.LatencyRecord.Windows,
				AlertWindow: slo.LatencyRecord.AlertWindow,
				AlertWait:   slo.LatencyRecord.AlertWait,
				BurnRate:    slo.ErrorRateRecord.BurnRate,
			})
			if err != nil {
				log.Panicf("Could not generate alert, err: %s", err.Error())
			}
			alertRules = append(alertRules, ruleNodes(latencyRules)...)
		}
	}

	for _, rule := range alertRules {
		slo.fillMetadata(&rule)
	}

	if disableTicket {
		var alertRulesWithoutTicket []rulefmt.RuleNode

		for _, rule := range alertRules {
			if rule.Labels["severity"] != "ticket" {
				alertRulesWithoutTicket = append(alertRulesWithoutTicket, rule)
			}
		}

		return alertRulesWithoutTicket
	}

	return alertRules
}

func (slo *SLO) fillMetadata(rule *rulefmt.RuleNode) {
	for label, value := range slo.Labels {
		rule.Labels[label] = value
	}

	for label, value := range slo.Annotations {
		rule.Annotations[label] = value
	}
}

func (slo *SLO) GenerateGroupRules(sloClass *Class, disableTicket bool) []rulefmt.RuleGroup {
	var rules []rulefmt.RuleGroup

	objectives := slo.Objectives
	if sloClass != nil {
		objectives = sloClass.Objectives
	}
	latencyBuckets := objectives.LatencyBuckets()
	if len(slo.LatencyRecord.Buckets) > 0 {
		latencyBuckets = slo.LatencyRecord.Buckets
	}

	for _, sample := range samples.DefaultSamples {

		interval, err := model.ParseDuration(sample.Interval)
		if err != nil {
			log.Fatal(err)
		}
		ruleGroup := rulefmt.RuleGroup{
			Name:     fmt.Sprintf("slo:%s:%s", slo.Name, sample.Name),
			Interval: interval,
			Rules:    []rulefmt.RuleNode{},
		}

		for _, bucket := range sample.Buckets {
			if disableTicket && samples.IsTicketSample(bucket) {
				continue
			}

			ruleGroup.Rules = append(ruleGroup.Rules, slo.generateRules(bucket, latencyBuckets)...)
		}

		if len(ruleGroup.Rules) > 0 {
			rules = append(rules, ruleGroup)
		}
	}

	return rules
}

func (slo *SLO) labels() map[string]string {
	labels := make(map[string]string)
	if !slo.HonorLabels {
		labels["service"] = slo.Name
	}
	for key, value := range slo.Labels {
		labels[key] = value
	}
	return labels
}

func (slo *SLO) generateRules(bucket string, latencyBuckets []string) []rulefmt.RuleNode {
	var rules []rulefmt.RuleNode
	if slo.TrafficRateRecord.Expr != "" {
		trafficRateRecord := rulefmt.RuleNode{
			Labels: slo.labels(),
		}

		trafficRateRecord.Record.SetString(fmt.Sprintf("slo:service_traffic:ratio_rate_%s", bucket))
		trafficRateRecord.Expr.SetString(slo.TrafficRateRecord.ComputeExpr(bucket, ""))

		rules = append(rules, trafficRateRecord)
	}

	if slo.ErrorRateRecord.Expr != "" {
		errorRateRecord := rulefmt.RuleNode{
			Labels: slo.labels(),
		}

		errorRateRecord.Record.SetString(fmt.Sprintf("slo:service_errors_total:ratio_rate_%s", bucket))
		errorRateRecord.Expr.SetString(slo.ErrorRateRecord.ComputeExpr(bucket, ""))
		rules = append(rules, errorRateRecord)
	}

	if slo.LatencyQuantileRecord.Expr != "" {
		for _, quantile := range quantiles {
			latencyQuantileRecord := rulefmt.RuleNode{
				Labels: slo.labels(),
			}

			latencyQuantileRecord.Record.SetString(fmt.Sprintf("slo:service_latency:%s_%s", quantile.name, bucket))
			latencyQuantileRecord.Expr.SetString(slo.LatencyQuantileRecord.ComputeQuantile(bucket, quantile.quantile))
			rules = append(rules, latencyQuantileRecord)
		}
	}

	if slo.LatencyRecord.Expr != "" {
		for _, latencyBucket := range latencyBuckets {
			latencyRateRecord := rulefmt.RuleNode{
				Labels: slo.labels(),
			}
			latencyRateRecord.Record.SetString("slo:service_latency:ratio_rate_" + bucket)
			latencyRateRecord.Expr.SetString(slo.LatencyRecord.ComputeExpr(bucket, latencyBucket))
			latencyRateRecord.Labels["le"] = latencyBucket

			rules = append(rules, latencyRateRecord)
		}
	}

	return rules
}

func ruleNodes(origin []rulefmt.Rule) []rulefmt.RuleNode {
	result := make([]rulefmt.RuleNode, len(origin))

	for i := range origin {
		result[i] = ruleNode(origin[i])
	}

	return result
}

func ruleNode(origin rulefmt.Rule) rulefmt.RuleNode {
	result := rulefmt.RuleNode{}

	if origin.Alert != "" {
		result.Alert.SetString(origin.Alert)
	}

	if origin.Record != "" {
		result.Record.SetString(origin.Record)
	}

	if origin.Expr != "" {
		result.Expr.SetString(origin.Expr)
	}

	result.For = origin.For
	result.Labels = origin.Labels
	result.Annotations = origin.Annotations

	return result
}
