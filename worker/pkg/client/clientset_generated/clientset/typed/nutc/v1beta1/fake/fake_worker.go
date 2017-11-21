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
package fake

import (
	v1beta1 "github.com/kairen/worker/pkg/apis/nutc/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWorkers implements WorkerInterface
type FakeWorkers struct {
	Fake *FakeNutcV1beta1
	ns   string
}

var workersResource = schema.GroupVersionResource{Group: "nutc.worker.com", Version: "v1beta1", Resource: "workers"}

var workersKind = schema.GroupVersionKind{Group: "nutc.worker.com", Version: "v1beta1", Kind: "Worker"}

// Get takes name of the worker, and returns the corresponding worker object, and an error if there is any.
func (c *FakeWorkers) Get(name string, options v1.GetOptions) (result *v1beta1.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(workersResource, c.ns, name), &v1beta1.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Worker), err
}

// List takes label and field selectors, and returns the list of Workers that match those selectors.
func (c *FakeWorkers) List(opts v1.ListOptions) (result *v1beta1.WorkerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(workersResource, workersKind, c.ns, opts), &v1beta1.WorkerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.WorkerList{}
	for _, item := range obj.(*v1beta1.WorkerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested workers.
func (c *FakeWorkers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(workersResource, c.ns, opts))

}

// Create takes the representation of a worker and creates it.  Returns the server's representation of the worker, and an error, if there is any.
func (c *FakeWorkers) Create(worker *v1beta1.Worker) (result *v1beta1.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(workersResource, c.ns, worker), &v1beta1.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Worker), err
}

// Update takes the representation of a worker and updates it. Returns the server's representation of the worker, and an error, if there is any.
func (c *FakeWorkers) Update(worker *v1beta1.Worker) (result *v1beta1.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(workersResource, c.ns, worker), &v1beta1.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Worker), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkers) UpdateStatus(worker *v1beta1.Worker) (*v1beta1.Worker, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(workersResource, "status", c.ns, worker), &v1beta1.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Worker), err
}

// Delete takes name of the worker and deletes it. Returns an error if one occurs.
func (c *FakeWorkers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(workersResource, c.ns, name), &v1beta1.Worker{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(workersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.WorkerList{})
	return err
}

// Patch applies the patch and returns the patched worker.
func (c *FakeWorkers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(workersResource, c.ns, name, data, subresources...), &v1beta1.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.Worker), err
}
