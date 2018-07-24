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

package operator

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
	opkit "github.com/inwinstack/operator-kit"
	exampleclientset "github.com/kairen/simple-operator/pkg/client/clientset/versioned/typed/example/v1alpha1"
	employee "github.com/kairen/simple-operator/pkg/operator/employee"
	"k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Option is the options of operator
type Option struct {
	Endpoint string
	Token    string
}

type Operator struct {
	ctx        *opkit.Context
	controller *employee.Controller
	resources  []opkit.CustomResource
	opts       Option
}

func NewMainOperator(opts Option) *Operator {
	schemes := []opkit.CustomResource{employee.Resource}
	o := &Operator{
		resources: schemes,
		opts:      opts,
	}
	return o
}

// Initialize init the instance resources.
func (o *Operator) Initialize() error {
	glog.V(2).Info("initialize the operator resources.")

	ctx, clientset, err := o.initContextAndClient()
	if err != nil {
		return err
	}
	o.controller = employee.NewController(ctx, clientset)
	o.ctx = ctx
	return nil
}

func (o *Operator) initContextAndClient() (*opkit.Context, exampleclientset.ExampleV1alpha1Interface, error) {
	glog.V(2).Info("initialize the operator context and client.")

	config, err := o.getRestConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Kubernetes config. %+v", err)
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get Kubernetes client. %+v", err)
	}

	ec, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes API extension clientset. %+v", err)
	}

	ic, err := exampleclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create IKM clientset. %+v", err)
	}

	ctx := &opkit.Context{
		Clientset:             c,
		APIExtensionClientset: ec,
		Interval:              OperatorInterval,
		Timeout:               OperatorTimeout,
	}
	return ctx, ic, nil
}

func (o *Operator) isDevMode() bool {
	if o.opts.Endpoint != "" && o.opts.Token != "" {
		return true
	}
	return false
}

func (o *Operator) getRestConfig() (*rest.Config, error) {
	dev := o.isDevMode()
	if !dev {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	config := &rest.Config{
		Host:            o.opts.Endpoint,
		BearerToken:     string(o.opts.Token),
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	return config, nil
}

func (o *Operator) Run() error {
	for {
		err := o.initResources()
		if err == nil {
			break
		}
		glog.Errorf("failed to init resources. %+v. retrying...", err)
		<-time.After(InitRetryDelay)
	}

	signalChan := make(chan os.Signal, 1)
	stopChan := make(chan struct{})
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	o.controller.StartWatch(v1.NamespaceAll, stopChan)
	for {
		select {
		case <-signalChan:
			glog.Infof("shutdown signal received, exiting...")
			close(stopChan)
			return nil
		}
	}
}

func (o *Operator) initResources() error {
	glog.V(2).Info("initialize the CRD resources.")

	ctx := opkit.Context{
		Clientset:             o.ctx.Clientset,
		APIExtensionClientset: o.ctx.APIExtensionClientset,
		Interval:              OperatorInterval,
		Timeout:               OperatorTimeout,
	}
	err := opkit.CreateCustomResources(ctx, o.resources)
	if err != nil {
		return fmt.Errorf("failed to create custom resource. %+v", err)
	}
	return nil
}
