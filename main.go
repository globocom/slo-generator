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
		sloPath       = ""
		classesPath   = ""
		ruleOutput    = ""
		disableTicket = false
	)
	flag.StringVar(&sloPath, "slo.path", "", "A YML file describing SLOs")
	flag.StringVar(&classesPath, "classes.path", "", "A YML file describing SLOs classes (optional)")
	flag.StringVar(&ruleOutput, "rule.output", "", "Output to describe a prometheus rules")
	flag.BoolVar(&disableTicket, "disable.ticket", false, "Disable generation of alerts of kind ticket")

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

	classesDefinition, err := readClassesDefinition(classesPath)
	if err != nil {
		log.Fatal(err)
	}

	ruleGroups := &rulefmt.RuleGroups{
		Groups: []rulefmt.RuleGroup{},
	}

	for _, slo := range spec.SLOS {
		// try to use any slo class found
		sloClass, err := classesDefinition.FindClass(slo.Class)
		if err != nil {
			log.Fatalf("Could not compile SLO: %q, err: %q", slo.Name, err.Error())
		}

		ruleGroups.Groups = append(ruleGroups.Groups, slo.GenerateGroupRules(sloClass, disableTicket)...)
		ruleGroups.Groups = append(ruleGroups.Groups, rulefmt.RuleGroup{
			Name:  "slo:" + slo.Name + ":alert",
			Rules: slo.GenerateAlertRules(sloClass, disableTicket),
		})
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

// readClassesDefinition read SLO classes from filesystem
func readClassesDefinition(classesPath string) (*slo.ClassesDefinition, error) {
	classesDefinition := slo.ClassesDefinition{
		Classes: []slo.Class{},
	}
	if classesPath != "" {
		f, err := os.Open(classesPath)
		if err != nil {
			return nil, err
		}
		err = yaml.NewDecoder(f).Decode(&classesDefinition)
		if err != nil {
			return nil, err
		}
	}

	return &classesDefinition, nil
}
