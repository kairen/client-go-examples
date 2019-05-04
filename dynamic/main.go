package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	filePath := flag.String("f", "", "path to the file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	f, err := os.Open(*filePath)
	if err != nil {
		klog.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewYAMLOrJSONDecoder(f, 4096)
	discoveryClient := clientset.Discovery()
	apigroups, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		klog.Fatal(err)
	}

	mapper := restmapper.NewDiscoveryRESTMapper(apigroups)
	for {
		ext := runtime.RawExtension{}
		if err := decoder.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			klog.Fatal(err)
		}

		versions := &runtime.VersionedObjects{}
		_, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, versions)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			klog.Fatal(err)
		}

		restconfig := config
		restconfig.GroupVersion = &schema.GroupVersion{
			Group:   mapping.GroupVersionKind.Group,
			Version: mapping.GroupVersionKind.Version,
		}

		var unstruct unstructured.Unstructured
		var blob interface{}
		unstruct.Object = make(map[string]interface{})
		if err := json.Unmarshal(ext.Raw, &blob); err != nil {
			klog.Fatal(err)
		}

		namespace := "default"
		unstruct.Object = blob.(map[string]interface{})
		if md, ok := unstruct.Object["metadata"]; ok {
			metadata := md.(map[string]interface{})
			if internalns, ok := metadata["namespace"]; ok {
				namespace = internalns.(string)
			}
		}

		if _, err := dynamicClient.Resource(mapping.Resource).Namespace(namespace).Create(&unstruct, metav1.CreateOptions{}); err != nil {
			klog.Fatal(err)
		}
		fmt.Printf("%s \"%s\" created.\n", mapping.GroupVersionKind.Kind, unstruct.GetName())
	}
}
