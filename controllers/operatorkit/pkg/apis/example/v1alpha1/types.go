/*
Copyright © 2018 Kyle Bai(kyle.b@inwinstack.com)

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Employee struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   EmployeeSpec   `json:"spec"`
	Status EmployeeStatus `json:"status"`
}

type EmployeeSpec struct {
	TaskName string `json:"taskName"`
	Threads  *int32 `json:"threads"`
}

type EmployeeLiverState string

const (
	LiverHealth    EmployeeLiverState = "Health"
	LiverFibrosis  EmployeeLiverState = "Fibrosis"
	LiverCirrhosis EmployeeLiverState = "Cirrhosis"
	LiverDead      EmployeeLiverState = "Dead"
)

type EmployeeStatus struct {
	AvailableThreads int32              `json:"availableThreads"`
	LiverState       EmployeeLiverState `json:"liverState"`
	Message          string             `json:"message"`
	LastLiveTime     metav1.Time        `json:"lastLiveTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type EmployeeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Employee `json:"items"`
}
