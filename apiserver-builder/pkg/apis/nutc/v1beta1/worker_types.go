/*
Copyright 2017 The Kubernetes Authors.

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

package v1beta1

import (
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/kairen/worker/pkg/apis/nutc"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +k8s:openapi-gen=true
// +resource:path=workers,strategy=WorkerStrategy
// Worker defines taks by assigner
type Worker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkerSpec   `json:"spec,omitempty"`
	Status WorkerStatus `json:"status,omitempty"`
}

// WorkerSpec defines the desired state of Worker
type WorkerSpec struct {
	Name     string   `json:"name,omitempty"`
	Assigner string   `json:"assigner,omitempty"`
	Tasks    []string `json:"tasks,omitempty"`
	WorkDay  int64    `json:"workDay,omitempty"`
}

// WorkerStatus defines the observed state of Worker
type WorkerStatus struct {
	State *string `json:"state,omitempty"`
}

// Validate checks that an instance of Worker is well formed
func (WorkerStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*nutc.Worker)
	log.Printf("Validating fields for Worker %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default Worker field values
func (WorkerSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*Worker)
	// set default field values here
	log.Printf("Defaulting fields for Worker %s\n", obj.Name)
	state := checkState(obj)
	obj.Status.State = &state
}

func checkState(obj *Worker) string {
	if obj.Spec.WorkDay <= 2 && len(obj.Spec.Tasks) >= 3 {
		return "Dead"
	} else if obj.Spec.WorkDay > 2 && len(obj.Spec.Tasks) < 2 {
		return "Easy"
	}
	return "Should not be a problem"
}
