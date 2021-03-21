/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2

import (
	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SLOSpec defines the desired state of SLO
type SLOSpec struct {
	Service    string     `json:"service,omitempty"` // Name of SLI or output of a multiService SLI
	Objectives Objectives `json:"objectives,omitempty"`
}

type Objectives struct {
	Availability float64         `json:"availability,omitempty"`
	Latency      []LatencyTarget `json:"latency,omitempty"`
	Window       model.Duration  `json:"window,omitempty"`
}

type Alert struct {
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Method      string            `json:"method,omitempty"`
	Window      string            `json:"window,omitempty"`
	BurnRate    float64           `json:"burnRate,omitempty"`
	Wait        string            `json:"wait,omitempty"`
}

type LatencyTarget struct {
	LE     string  `yaml:"le"`
	Target float64 `yaml:"target"`
}

// SLOStatus defines the observed state of SLO
type SLOStatus struct {
	UpdatedRevision int `json:"updatedRevision,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SLO is the Schema for the sloes API
type SLO struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SLOSpec   `json:"spec,omitempty"`
	Status SLOStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SLOList contains a list of SLO
type SLOList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SLO `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SLO{}, &SLOList{})
}
