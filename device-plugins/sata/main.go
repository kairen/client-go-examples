package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: sdp -stderrthreshold=[INFO|WARN|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func main() {
	glog.V(2).Infof("Starting SATA Device Plugin.")
	glog.V(2).Infof("Starting file system watcher.")
	watcher, err := newFSWatcher(pluginapi.DevicePluginPath)

	if err != nil {
		glog.V(2).Infof("Failed to created file system watcher.")
		os.Exit(1)
	}
	defer watcher.Close()

	glog.V(2).Infof("Starting operating system watcher.")
	sigs := newOSWatcher(syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	restart := true
	var devicePlugin *KairenDevicePlugin

L:
	for {
		if restart {
			if devicePlugin != nil {
				devicePlugin.Stop()
			}

			devicePlugin, err = NewKairenDevicePlugin()
			if err != nil {
				glog.V(2).Infof("Failed to new device plugin. %s", err)
				os.Exit(1)
			}

			if err := devicePlugin.Serve(); err != nil {
				glog.V(2).Infof("Could not contact Kubelet, retrying. Did you enable the device plugin feature gate?")
			} else {
				restart = false
			}
		}

		select {
		case event := <-watcher.Events:
			if event.Name == pluginapi.KubeletSocket && event.Op&fsnotify.Create == fsnotify.Create {
				glog.V(2).Infof("inotify: %s created, restarting.", pluginapi.KubeletSocket)
				restart = true
			}

			if event.Name == serverSock && event.Op&fsnotify.Remove == fsnotify.Remove {
				glog.V(2).Infof("inotify: %s deleted, restarting.", serverSock)
				restart = true
			}

		case err := <-watcher.Errors:
			glog.V(2).Infof("inotify: %s", err)
		case s := <-sigs:
			switch s {
			case syscall.SIGHUP:
				glog.V(2).Infof("Received SIGHUP, restarting.")
				restart = true
			default:
				glog.V(2).Infof("Received signal \"%v\", shutting down.", s)
				devicePlugin.Stop()
				break L
			}
		}
	}
	glog.Flush()
}
