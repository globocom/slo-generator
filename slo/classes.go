package slo

// Class representing a template of objectives
// this is important to achieve a scalabe SLO policies
// read more at: https://landing.google.com/sre/workbook/chapters/alerting-on-slos/#alerting_at_scale
type Class struct {
	Name       string     `yaml:"name"`
	Objectives Objectives `yaml:"objectives"`
}

type ClassesDefinition struct {
	Classes []Class `yaml:"classes"`
}

func (c *ClassesDefinition) FindClass(name string) *Class {
	for _, class := range c.Classes {
		if class.Name == name {
			return &class
		}
	}
	return nil
}
