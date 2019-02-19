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
package internalversion

import (
	nutc "github.com/kairen/worker/pkg/apis/nutc"
	scheme "github.com/kairen/worker/pkg/client/clientset_generated/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WorkersGetter has a method to return a WorkerInterface.
// A group's client should implement this interface.
type WorkersGetter interface {
	Workers(namespace string) WorkerInterface
}

// WorkerInterface has methods to work with Worker resources.
type WorkerInterface interface {
	Create(*nutc.Worker) (*nutc.Worker, error)
	Update(*nutc.Worker) (*nutc.Worker, error)
	UpdateStatus(*nutc.Worker) (*nutc.Worker, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*nutc.Worker, error)
	List(opts v1.ListOptions) (*nutc.WorkerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *nutc.Worker, err error)
	WorkerExpansion
}

// workers implements WorkerInterface
type workers struct {
	client rest.Interface
	ns     string
}

// newWorkers returns a Workers
func newWorkers(c *NutcClient, namespace string) *workers {
	return &workers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the worker, and returns the corresponding worker object, and an error if there is any.
func (c *workers) Get(name string, options v1.GetOptions) (result *nutc.Worker, err error) {
	result = &nutc.Worker{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Workers that match those selectors.
func (c *workers) List(opts v1.ListOptions) (result *nutc.WorkerList, err error) {
	result = &nutc.WorkerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested workers.
func (c *workers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("workers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a worker and creates it.  Returns the server's representation of the worker, and an error, if there is any.
func (c *workers) Create(worker *nutc.Worker) (result *nutc.Worker, err error) {
	result = &nutc.Worker{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("workers").
		Body(worker).
		Do().
		Into(result)
	return
}

// Update takes the representation of a worker and updates it. Returns the server's representation of the worker, and an error, if there is any.
func (c *workers) Update(worker *nutc.Worker) (result *nutc.Worker, err error) {
	result = &nutc.Worker{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("workers").
		Name(worker.Name).
		Body(worker).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *workers) UpdateStatus(worker *nutc.Worker) (result *nutc.Worker, err error) {
	result = &nutc.Worker{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("workers").
		Name(worker.Name).
		SubResource("status").
		Body(worker).
		Do().
		Into(result)
	return
}

// Delete takes name of the worker and deletes it. Returns an error if one occurs.
func (c *workers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *workers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched worker.
func (c *workers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *nutc.Worker, err error) {
	result = &nutc.Worker{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("workers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
