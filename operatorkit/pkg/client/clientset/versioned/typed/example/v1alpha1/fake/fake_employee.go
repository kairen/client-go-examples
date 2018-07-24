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
package fake

import (
	v1alpha1 "github.com/kairen/simple-operator/pkg/apis/example/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEmployees implements EmployeeInterface
type FakeEmployees struct {
	Fake *FakeExampleV1alpha1
	ns   string
}

var employeesResource = schema.GroupVersionResource{Group: "example.io", Version: "v1alpha1", Resource: "employees"}

var employeesKind = schema.GroupVersionKind{Group: "example.io", Version: "v1alpha1", Kind: "Employee"}

// Get takes name of the employee, and returns the corresponding employee object, and an error if there is any.
func (c *FakeEmployees) Get(name string, options v1.GetOptions) (result *v1alpha1.Employee, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(employeesResource, c.ns, name), &v1alpha1.Employee{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Employee), err
}

// List takes label and field selectors, and returns the list of Employees that match those selectors.
func (c *FakeEmployees) List(opts v1.ListOptions) (result *v1alpha1.EmployeeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(employeesResource, employeesKind, c.ns, opts), &v1alpha1.EmployeeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EmployeeList{}
	for _, item := range obj.(*v1alpha1.EmployeeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested employees.
func (c *FakeEmployees) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(employeesResource, c.ns, opts))

}

// Create takes the representation of a employee and creates it.  Returns the server's representation of the employee, and an error, if there is any.
func (c *FakeEmployees) Create(employee *v1alpha1.Employee) (result *v1alpha1.Employee, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(employeesResource, c.ns, employee), &v1alpha1.Employee{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Employee), err
}

// Update takes the representation of a employee and updates it. Returns the server's representation of the employee, and an error, if there is any.
func (c *FakeEmployees) Update(employee *v1alpha1.Employee) (result *v1alpha1.Employee, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(employeesResource, c.ns, employee), &v1alpha1.Employee{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Employee), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEmployees) UpdateStatus(employee *v1alpha1.Employee) (*v1alpha1.Employee, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(employeesResource, "status", c.ns, employee), &v1alpha1.Employee{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Employee), err
}

// Delete takes name of the employee and deletes it. Returns an error if one occurs.
func (c *FakeEmployees) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(employeesResource, c.ns, name), &v1alpha1.Employee{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEmployees) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(employeesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.EmployeeList{})
	return err
}

// Patch applies the patch and returns the patched employee.
func (c *FakeEmployees) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Employee, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(employeesResource, c.ns, name, data, subresources...), &v1alpha1.Employee{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Employee), err
}
