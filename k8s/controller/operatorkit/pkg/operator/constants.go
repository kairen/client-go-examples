/*
Copyright © 2018 inwinSTACK.inc

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

package operator

import "time"

const (
	// PodNameEnvVar is the env variable for getting the pod name via downward api
	PodNameEnvVar = "POD_NAME"

	// PodNamespaceEnvVar is the env variable for getting the pod namespace via downward api
	PodNamespaceEnvVar = "POD_NAMESPACE"

	// NodeNameEnvVar is the env variable for getting the node via downward api
	NodeNameEnvVar = "NODE_NAME"

	// InitRetryDelay is the operator init retry delay time
	InitRetryDelay = 10 * time.Second

	// OperatorInterval is the operator connect interval time
	OperatorInterval = 500 * time.Millisecond

	// OperatorTimeout is the operator connect timeout time
	OperatorTimeout = 60 * time.Second
)
