package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	ghodssYaml "github.com/ghodss/yaml"
	"github.com/globocom/slo-generator/kubernetes"
	"github.com/globocom/slo-generator/slo"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/prometheus/prometheus/pkg/rulefmt"
	yaml "gopkg.in/yaml.v3"
)

func main() {
	var (
		sloPath       = ""
		classesPath   = ""
		ruleOutput    = ""
		k8sLabels     = ""
		disableTicket = false
		k8s           = false
	)
	flag.StringVar(&sloPath, "slo.path", "", "A YML file describing SLOs")
	flag.StringVar(&classesPath, "classes.path", "", "A YML file describing SLOs classes (optional)")
	flag.StringVar(&ruleOutput, "rule.output", "", "Output to describe a prometheus rules")
	flag.BoolVar(&disableTicket, "disable.ticket", false, "Disable generation of alerts of kind ticket")
	flag.BoolVar(&k8s, "kubernetes", false, "Generates prometheus-operator YAML")
	flag.StringVar(&k8sLabels, "kubernetes-labels", "", "Add some labels in generated resource")

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

	if len(spec.Classes) > 0 && len(classesDefinition.Classes) > 0 {
		log.Fatal("you can not define classes in slo and classes files")
	} else if len(classesDefinition.Classes) > 0 {
		spec.Classes = classesDefinition.Classes
	}

	ruleGroups := &rulefmt.RuleGroups{
		Groups: []rulefmt.RuleGroup{},
	}

	if k8s {
		labels := map[string]string{}
		if k8sLabels != "" {
			labels, err = parseLabels(k8sLabels)
			if err != nil {
				log.Fatal(err)
			}
		}

		manifests := []monitoringv1.PrometheusRule{}
		for _, slo := range spec.SLOS {
			// try to use any slo class found
			sloClass, err := spec.Classes.FindClass(slo.Class)
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
			if manifest.Labels == nil {
				manifest.Labels = map[string]string{}
			}
			for key, value := range labels {
				manifest.Labels[key] = value
			}

			b, err := ghodssYaml.Marshal(manifest)

			if err != nil {
				log.Fatal(err)
			}
			if i > 0 {
				_, err = output.Write([]byte("---\n"))
				if err != nil {
					log.Fatal(err)
				}
			}
			_, err = output.Write(b)
			if err != nil {
				log.Fatal(err)
			}
		}
		if ruleOutput != "" {
			log.Printf("generated a kubernetes manifest record in %q", ruleOutput)
		}
		return
	}

	for _, slo := range spec.SLOS {
		// try to use any slo class found
		sloClass, err := spec.Classes.FindClass(slo.Class)
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

func parseLabels(labels string) (map[string]string, error) {
	result := map[string]string{}
	for _, part := range strings.Split(labels, ",") {
		parts := strings.Split(part, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid label %q", part)
		}
		result[parts[0]] = parts[1]
	}

	return result, nil
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
