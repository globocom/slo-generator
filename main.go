package main

import (
	"flag"
	"log"
	"os"
	"strings"

	algorithms "github.com/globocom/slo-generator/algorithms"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v2"
)

type Sample struct {
	Name     string
	Interval string
	Buckets  []string
}

var defaultSamples = []Sample{
	{
		Name:     "short",
		Interval: "30s",
		Buckets:  []string{"5m", "30m", "1h"},
	},
	{
		Name:     "medium",
		Interval: "2m",
		Buckets:  []string{"2h", "6h"},
	},
	{
		Name:     "daily",
		Interval: "5m",
		Buckets:  []string{"1d", "3d"},
	},
}

type SLOSpec struct {
	SLOS []SLO
}

type ExprBlock struct {
	Expr string `yaml:"expr"`
}

func (block *ExprBlock) ComputeExpr(window string) string {
	replacer := strings.NewReplacer("$window", window)
	return replacer.Replace(block.Expr)
}

type SLO struct {
	Name                         string             `yaml:"name"`
	Algorithm                    string             `yaml:"algorithm"`
	AvailabilityObjectivePercent float64            `yaml:"availabilityObjectivePercent"`
	LatencyObjectiveBuckets      map[float64]string `yaml:"latencyObjectiveBuckets"`
	ErrorRateRecord              ExprBlock          `yaml:"errorRateRecord"`
	LatencyRecord                ExprBlock          `yaml:"latencyRecord"`
	Annotations                  map[string]string  `yaml:"annotations"`
}

func (slo SLO) GenerateGroupRules() []rulefmt.RuleGroup {
	rules := []rulefmt.RuleGroup{}

	for _, sample := range defaultSamples {
		interval, err := model.ParseDuration(sample.Interval)
		if err != nil {
			log.Fatal(err)
		}
		ruleGroup := rulefmt.RuleGroup{
			Name:     "slo:" + slo.Name + "_" + sample.Name,
			Interval: interval,
			Rules:    []rulefmt.Rule{},
		}

		for _, bucket := range sample.Buckets {
			errorRateRecord := rulefmt.Rule{
				Record: "slo:service_errors_total:ratio_rate_" + bucket,
				Expr:   slo.ErrorRateRecord.ComputeExpr(bucket),
				Labels: map[string]string{
					"service": slo.Name,
				},
			}

			ruleGroup.Rules = append(ruleGroup.Rules, errorRateRecord)
		}

		rules = append(rules, ruleGroup)
	}

	// alerting
	alertingGroup := rulefmt.RuleGroup{
		Name:  "slo:" + slo.Name + "_alert",
		Rules: []rulefmt.Rule{},
	}

	// alerting page
	sloPageRecord := rulefmt.Rule{
		Alert: "slo:" + slo.Name + ".errors.page",
		Expr: algorithms.MultiBurnRateForPage(
			"slo:service_errors_total",
			labels.New(labels.Label{"service", slo.Name}),
			"<", (1 - slo.AvailabilityObjectivePercent/100),
		),
		Annotations: slo.Annotations,
	}

	alertingGroup.Rules = append(alertingGroup.Rules, sloPageRecord)

	// alerting ticket
	sloTicketRecord := rulefmt.Rule{
		Alert: "slo:" + slo.Name + ".errors.ticket",
		Expr: algorithms.MultiBurnRateForTicket(
			"slo:service_errors_total",
			labels.New(labels.Label{"service", slo.Name}),
			"<", (1 - slo.AvailabilityObjectivePercent/100),
		),
		Annotations: slo.Annotations,
	}

	alertingGroup.Rules = append(alertingGroup.Rules, sloTicketRecord)

	rules = append(rules, alertingGroup)

	return rules
}

func main() {
	var (
		sloPath    = ""
		ruleOutput = ""
	)
	flag.StringVar(&sloPath, "slo.path", "", "A YML file describing SLOs")
	flag.StringVar(&ruleOutput, "rule.output", "", "Output to describe a prometheus rules")

	flag.Parse()

	if sloPath == "" {
		log.Fatal("slo.path is a required param")
	}

	if ruleOutput == "" {
		log.Fatal("rule.output is a required param")
	}

	f, err := os.Open(sloPath)
	if err != nil {
		log.Fatal(err)
	}

	spec := &SLOSpec{}
	err = yaml.NewDecoder(f).Decode(spec)
	if err != nil {
		log.Fatal(err)
	}

	ruleGroups := &rulefmt.RuleGroups{
		Groups: []rulefmt.RuleGroup{},
	}

	for _, slo := range spec.SLOS {
		ruleGroups.Groups = append(ruleGroups.Groups, slo.GenerateGroupRules()...)
	}

	targetFile, err := os.Create(ruleOutput)
	if err != nil {
		log.Fatal(err)
	}
	defer targetFile.Close()
	err = yaml.NewEncoder(targetFile).Encode(ruleGroups)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("generated a SLO record in %q", ruleOutput)
}
