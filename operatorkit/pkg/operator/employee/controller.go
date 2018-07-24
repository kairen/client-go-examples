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

package employee

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"

	opkit "github.com/inwinstack/operator-kit"
	examplev1alpha1 "github.com/kairen/simple-operator/pkg/apis/example/v1alpha1"
	exampleclientset "github.com/kairen/simple-operator/pkg/client/clientset/versioned/typed/example/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/tools/cache"
)

const (
	CustomResourceName       = "employee"
	CustomResourceNamePlural = "employees"
	createInterval           = 6 * time.Second
	createTimeout            = 5 * time.Minute
	namespace                = "default"
)

var Resource = opkit.CustomResource{
	Name:       CustomResourceName,
	Plural:     CustomResourceNamePlural,
	Group:      examplev1alpha1.CustomResourceGroup,
	Version:    examplev1alpha1.Version,
	Scope:      apiextensionsv1beta1.NamespaceScoped,
	Kind:       reflect.TypeOf(examplev1alpha1.Employee{}).Name(),
	ShortNames: []string{"emp"},
}

var msg = map[examplev1alpha1.EmployeeLiverState]string{
	examplev1alpha1.LiverHealth:    "這工作太簡單了!我非常健康!!",
	examplev1alpha1.LiverFibrosis:  "我的肝臟正在逐漸惡化與纖維化!!",
	examplev1alpha1.LiverCirrhosis: "我的肝臟已經硬化了!!",
	examplev1alpha1.LiverDead:      "我的肝臟已經損壞了!!",
}

type Controller struct {
	context   *opkit.Context
	clientset exampleclientset.ExampleV1alpha1Interface
	clock     clock.Clock
}

func NewController(context *opkit.Context, clientset exampleclientset.ExampleV1alpha1Interface) *Controller {
	return &Controller{
		context:   context,
		clientset: clientset,
		clock:     clock.RealClock{},
	}
}

func (c *Controller) StartWatch(namespace string, stopCh chan struct{}) error {
	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	glog.Infof("start watching resources in namespace %s", namespace)
	watcher := opkit.NewWatcher(Resource, namespace, resourceHandlers, c.clientset.RESTClient())
	go watcher.Watch(&examplev1alpha1.Employee{}, stopCh)
	return nil
}

func (c *Controller) onAdd(obj interface{}) {
	employee := obj.(*examplev1alpha1.Employee).DeepCopy()
	glog.Infof("%s resource onAdd.", employee.Name)

	deployment, _ := c.context.Clientset.AppsV1().Deployments(employee.Namespace).Create(newDeployment(employee))
	c.updateStatus(employee, deployment)
}

func (c *Controller) onUpdate(oldObj, newObj interface{}) {
	employee := oldObj.(*examplev1alpha1.Employee).DeepCopy()
	newEmployee := newObj.(*examplev1alpha1.Employee).DeepCopy()
	glog.Infof("%s resource onUpdate.", employee.Name)

	deployment, _ := c.context.Clientset.AppsV1().Deployments(employee.Namespace).Update(newDeployment(employee))
	c.updateStatus(newEmployee, deployment)
}

func (c *Controller) onDelete(obj interface{}) {
	employee := obj.(*examplev1alpha1.Employee).DeepCopy()
	glog.Infof("%s resource onDelete.", employee.Name)

	c.context.Clientset.AppsV1().Deployments(employee.Namespace).Delete(employee.Name, &metav1.DeleteOptions{})
}

func (c *Controller) updateStatus(employee *examplev1alpha1.Employee, deployment *appsv1.Deployment) error {
	employeeCopy := employee.DeepCopy()
	employeeCopy.Status.AvailableThreads = *employee.Spec.Threads

	r := employeeCopy.Status.AvailableThreads
	switch {
	case r <= 1:
		employeeCopy.Status.LiverState = examplev1alpha1.LiverHealth
	case r <= 3 && r > 1:
		employeeCopy.Status.LiverState = examplev1alpha1.LiverFibrosis
	case r > 3 && r <= 9:
		employeeCopy.Status.LiverState = examplev1alpha1.LiverCirrhosis
	case r > 9:
		employeeCopy.Status.LiverState = examplev1alpha1.LiverDead
	}

	employeeCopy.Status.Message = msg[employeeCopy.Status.LiverState]
	employeeCopy.Status.LastLiveTime = metav1.NewTime(c.clock.Now())
	if _, err := c.clientset.Employees(namespace).Update(employeeCopy); err != nil {
		return fmt.Errorf("failed to update employee %s status: %+v", employee.Namespace, err)
	}
	return nil
}

func newDeployment(employee *examplev1alpha1.Employee) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": employee.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      employee.Spec.TaskName,
			Namespace: employee.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(employee, schema.GroupVersionKind{
					Group:   examplev1alpha1.SchemeGroupVersion.Group,
					Version: examplev1alpha1.SchemeGroupVersion.Version,
					Kind:    "Employee",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: employee.Spec.Threads,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}
