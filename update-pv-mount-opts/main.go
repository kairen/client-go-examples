package main

import (
	goflag "flag"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	slice "github.com/thoas/go-funk"

	flag "github.com/spf13/pflag"
)

var (
	kubeconfig string
)

func parserFlags() {
	flag.StringVarP(&kubeconfig, "kubeconfig", "", "", "Absolute path to the kubeconfig file.")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()
}

func buildRestConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	defer glog.Flush()
	parserFlags()

	config, err := buildRestConfig(kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to get Kubernetes config. %+v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to get Kubernetes client. %+v", err)
	}

	pvList, err := client.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
	if err != nil {
		glog.Fatalf("Failed to get PVs. %+v", err)
	}

	for _, pv := range pvList.Items {
		if slice.ContainsString(pv.Spec.MountOptions, "vers=4.1") {
			glog.Infof("PV %s has mount options for NFS v4.1.", pv.Name)
			pv.Spec.MountOptions = slice.FilterString(pv.Spec.MountOptions, func(v string) bool {
				return v != "vers=4.1"
			})
			pv.Spec.MountOptions = append(pv.Spec.MountOptions, "vers=4.0")

			_, err := client.CoreV1().PersistentVolumes().Update(&pv)
			if err != nil {
				glog.Fatalf("Failed to update  %s PV. %+v", pv.Name, err)
			}
			glog.Infof("PV %s has been updated mount options for NFS v4.0.", pv.Name)
		}
	}
}
