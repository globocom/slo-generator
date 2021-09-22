package kubernetes

import (
	"github.com/globocom/slo-generator/slo"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Opts struct {
	SLO           slo.SLO
	Class         *slo.Class
	DisableTicket bool
}

func GenerateManifests(opt Opts) []monitoringv1.PrometheusRule {
	rules := []monitoringv1.PrometheusRule{}

	groups := opt.SLO.GenerateGroupRules(opt.Class, opt.DisableTicket)
	if len(groups) > 0 {
		rules = append(rules, monitoringv1.PrometheusRule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PrometheusRule",
				APIVersion: "monitoring.coreos.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "slis-" + opt.SLO.Name,
			},
			Spec: monitoringv1.PrometheusRuleSpec{
				Groups: kubernetizeRuleGroups(groups),
			},
		})
	}

	alertRules := opt.SLO.GenerateAlertRules(opt.Class, opt.DisableTicket)
	if len(alertRules) > 0 {
		rules = append(rules, monitoringv1.PrometheusRule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PrometheusRule",
				APIVersion: "monitoring.coreos.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "slos-alerts-" + opt.SLO.Name,
			},
			Spec: monitoringv1.PrometheusRuleSpec{
				Groups: kubernetizeRuleGroups([]rulefmt.RuleGroup{
					{
						Name:  "slo:" + opt.SLO.Name + ":alert",
						Rules: alertRules,
					},
				}),
			},
		})

	}

	return rules
}

func kubernetizeRuleGroups(groups []rulefmt.RuleGroup) []monitoringv1.RuleGroup {
	result := []monitoringv1.RuleGroup{}
	for _, group := range groups {
		groupInterval := group.Interval.String()
		if groupInterval == "0s" {
			groupInterval = ""
		}
		result = append(result, monitoringv1.RuleGroup{
			Name:     group.Name,
			Interval: groupInterval,
			Rules:    kubernetizeRules(group.Rules),
		})
	}
	return result
}

func kubernetizeRules(rules []rulefmt.RuleNode) []monitoringv1.Rule {
	result := []monitoringv1.Rule{}
	for _, rule := range rules {
		ruleFor := rule.For.String()
		if ruleFor == "0s" {
			ruleFor = ""
		}
		result = append(result, monitoringv1.Rule{
			Record:      rule.Record.Value,
			Alert:       rule.Alert.Value,
			Expr:        intstr.FromString(rule.Expr.Value),
			For:         ruleFor,
			Labels:      rule.Labels,
			Annotations: rule.Annotations,
		})
	}
	return result
}
