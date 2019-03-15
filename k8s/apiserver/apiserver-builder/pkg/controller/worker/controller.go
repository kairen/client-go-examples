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

package worker

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"github.com/kairen/worker/pkg/apis/nutc/v1beta1"
	listers "github.com/kairen/worker/pkg/client/listers_generated/nutc/v1beta1"
	"github.com/kairen/worker/pkg/controller/sharedinformers"
)

// +controller:group=nutc,version=v1beta1,kind=Worker,resource=workers
type WorkerControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about Worker
	lister listers.WorkerLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *WorkerControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing workers labels
	c.lister = arguments.GetSharedInformers().Factory.Nutc().V1beta1().Workers().Lister()
}

// Reconcile handles enqueued messages
func (c *WorkerControllerImpl) Reconcile(u *v1beta1.Worker) error {
	// Implement controller logic here
	log.Printf("Running reconcile Worker for %s\n", u.Name)
	return nil
}

func (c *WorkerControllerImpl) Get(namespace, name string) (*v1beta1.Worker, error) {
	return c.lister.Workers(namespace).Get(name)
}
