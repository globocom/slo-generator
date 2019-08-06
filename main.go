package main

import (
	"flag"
	"log"
	"os"

	"github.com/globocom/slo-generator/slo"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v2"
)

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

	spec := &slo.SLOSpec{}
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
