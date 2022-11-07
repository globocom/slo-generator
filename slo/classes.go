package slo

import "fmt"

// Class represents a template of objectives
// this is important to achieve scalable SLO policies
// read more at: https://landing.google.com/sre/workbook/chapters/alerting-on-slos/#alerting_at_scale
type Class struct {
	Name       string     `yaml:"name"`
	Objectives Objectives `yaml:"objectives"`
}

type ClassesDefinition struct {
	Classes Classes `yaml:"classes"`
}

type Classes []Class

// FindClass finds for a given name, if not found return an error
func (c Classes) FindClass(name string) (*Class, error) {
	if name == "" {
		return nil, nil
	}
	for _, class := range c {
		if class.Name == name {
			return &class, nil
		}
	}

	return nil, fmt.Errorf("SLO class %q is not found", name)
}
