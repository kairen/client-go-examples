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
	nutc "github.com/kairen/worker/pkg/apis/nutc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWorkers implements WorkerInterface
type FakeWorkers struct {
	Fake *FakeNutc
	ns   string
}

var workersResource = schema.GroupVersionResource{Group: "nutc", Version: "", Resource: "workers"}

var workersKind = schema.GroupVersionKind{Group: "nutc", Version: "", Kind: "Worker"}

// Get takes name of the worker, and returns the corresponding worker object, and an error if there is any.
func (c *FakeWorkers) Get(name string, options v1.GetOptions) (result *nutc.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(workersResource, c.ns, name), &nutc.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*nutc.Worker), err
}

// List takes label and field selectors, and returns the list of Workers that match those selectors.
func (c *FakeWorkers) List(opts v1.ListOptions) (result *nutc.WorkerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(workersResource, workersKind, c.ns, opts), &nutc.WorkerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &nutc.WorkerList{}
	for _, item := range obj.(*nutc.WorkerList).Items {
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
func (c *FakeWorkers) Create(worker *nutc.Worker) (result *nutc.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(workersResource, c.ns, worker), &nutc.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*nutc.Worker), err
}

// Update takes the representation of a worker and updates it. Returns the server's representation of the worker, and an error, if there is any.
func (c *FakeWorkers) Update(worker *nutc.Worker) (result *nutc.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(workersResource, c.ns, worker), &nutc.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*nutc.Worker), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkers) UpdateStatus(worker *nutc.Worker) (*nutc.Worker, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(workersResource, "status", c.ns, worker), &nutc.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*nutc.Worker), err
}

// Delete takes name of the worker and deletes it. Returns an error if one occurs.
func (c *FakeWorkers) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(workersResource, c.ns, name), &nutc.Worker{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(workersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &nutc.WorkerList{})
	return err
}

// Patch applies the patch and returns the patched worker.
func (c *FakeWorkers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *nutc.Worker, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(workersResource, c.ns, name, data, subresources...), &nutc.Worker{})

	if obj == nil {
		return nil, err
	}
	return obj.(*nutc.Worker), err
}
