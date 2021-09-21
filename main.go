package main

import (
	"flag"
	"io"
	"log"
	"os"

	ghodssYaml "github.com/ghodss/yaml"
	"github.com/globocom/slo-generator/kubernetes"
	"github.com/globocom/slo-generator/slo"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	var (
		sloPath       = ""
		classesPath   = ""
		ruleOutput    = ""
		disableTicket = false
		k8s           = false
	)
	flag.StringVar(&sloPath, "slo.path", "", "A YML file describing SLOs")
	flag.StringVar(&classesPath, "classes.path", "", "A YML file describing SLOs classes (optional)")
	flag.StringVar(&ruleOutput, "rule.output", "", "Output to describe a prometheus rules")
	flag.BoolVar(&disableTicket, "disable.ticket", false, "Disable generation of alerts of kind ticket")
	flag.BoolVar(&k8s, "kubernetes", false, "Generates prometheus-operator YAML")

	flag.Parse()

	if sloPath == "" {
		log.Fatal("slo.path is a required param")
	}

	var output io.Writer
	if ruleOutput == "" {
		output = os.Stdout
	} else {
		targetFile, err := os.Create(ruleOutput)
		if err != nil {
			log.Fatal(err)
		}
		defer targetFile.Close()
		output = targetFile
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

	if k8s {
		manifests := []monitoringv1.PrometheusRule{}
		for _, slo := range spec.SLOS {
			// try to use any slo class found
			sloClass, err := classesDefinition.FindClass(slo.Class)
			if err != nil {
				log.Fatalf("Could not compile SLO: %q, err: %q", slo.Name, err.Error())
			}

			manifests = append(manifests, kubernetes.GenerateManifests(kubernetes.Opts{
				SLO:           slo,
				Class:         sloClass,
				DisableTicket: disableTicket,
			})...)
		}

		for i, manifest := range manifests {
			b, err := ghodssYaml.Marshal(manifest)

			if err != nil {
				log.Fatal(err)
			}
			if i > 0 {
				output.Write([]byte("---\n"))
			}
			output.Write(b)
		}
		if ruleOutput != "" {
			log.Printf("generated a kubernetes manifest record in %q", ruleOutput)
		}
		return
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

	err = yaml.NewEncoder(output).Encode(ruleGroups)
	if err != nil {
		log.Fatal(err)
	}
	if ruleOutput != "" {
		log.Printf("generated a SLO record in %q", ruleOutput)
	}
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
