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
	v1alpha1 "github.com/kairen/simple-operator/pkg/apis/example/v1alpha1"
	scheme "github.com/kairen/simple-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EmployeesGetter has a method to return a EmployeeInterface.
// A group's client should implement this interface.
type EmployeesGetter interface {
	Employees(namespace string) EmployeeInterface
}

// EmployeeInterface has methods to work with Employee resources.
type EmployeeInterface interface {
	Create(*v1alpha1.Employee) (*v1alpha1.Employee, error)
	Update(*v1alpha1.Employee) (*v1alpha1.Employee, error)
	UpdateStatus(*v1alpha1.Employee) (*v1alpha1.Employee, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Employee, error)
	List(opts v1.ListOptions) (*v1alpha1.EmployeeList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Employee, err error)
	EmployeeExpansion
}

// employees implements EmployeeInterface
type employees struct {
	client rest.Interface
	ns     string
}

// newEmployees returns a Employees
func newEmployees(c *ExampleV1alpha1Client, namespace string) *employees {
	return &employees{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the employee, and returns the corresponding employee object, and an error if there is any.
func (c *employees) Get(name string, options v1.GetOptions) (result *v1alpha1.Employee, err error) {
	result = &v1alpha1.Employee{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("employees").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Employees that match those selectors.
func (c *employees) List(opts v1.ListOptions) (result *v1alpha1.EmployeeList, err error) {
	result = &v1alpha1.EmployeeList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("employees").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested employees.
func (c *employees) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("employees").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a employee and creates it.  Returns the server's representation of the employee, and an error, if there is any.
func (c *employees) Create(employee *v1alpha1.Employee) (result *v1alpha1.Employee, err error) {
	result = &v1alpha1.Employee{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("employees").
		Body(employee).
		Do().
		Into(result)
	return
}

// Update takes the representation of a employee and updates it. Returns the server's representation of the employee, and an error, if there is any.
func (c *employees) Update(employee *v1alpha1.Employee) (result *v1alpha1.Employee, err error) {
	result = &v1alpha1.Employee{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("employees").
		Name(employee.Name).
		Body(employee).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *employees) UpdateStatus(employee *v1alpha1.Employee) (result *v1alpha1.Employee, err error) {
	result = &v1alpha1.Employee{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("employees").
		Name(employee.Name).
		SubResource("status").
		Body(employee).
		Do().
		Into(result)
	return
}

// Delete takes name of the employee and deletes it. Returns an error if one occurs.
func (c *employees) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("employees").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *employees) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("employees").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched employee.
func (c *employees) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Employee, err error) {
	result = &v1alpha1.Employee{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("employees").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
