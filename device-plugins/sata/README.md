[![Build Status](https://travis-ci.org/kairen/simple-device-plugin.svg?branch=master)](https://travis-ci.org/kairen/simple-device-plugin) [![Docker Build Statu](https://img.shields.io/docker/build/kairen/simple-device-plugin.svg)](https://hub.docker.com/r/kairen/simple-device-plugin/) [![codecov](https://codecov.io/gh/kairen/simple-device-plugin/branch/master/graph/badge.svg)](https://codecov.io/gh/kairen/simple-device-plugin)
# Simple Kubernetes Device Plugins
Learning how to implement the Kubernetes device plugins. This device plugin will automatically maps the SATA device according to your container SATA requirement.

## Prerequisites
The list of prerequisites for running the SATA device plugin is described below:
* Kubernetes version = 1.10.x
* The `DevicePlugins` feature gate enabled

## Quick Start
To build image:
```sh
$ make build_image
```

To install the SATA device plugins:
```sh
$ kubectl create -f https://raw.githubusercontent.com/kairen/simple-device-plugin/master/artifacts/simple-device-plugin.yml
$ kubectl -n kube-system get po -l name=device-plugin
NAME                            READY     STATUS    RESTARTS   AGE
simple-device-plugin-ds-jlj8k   1/1       Running   0          38s
simple-device-plugin-ds-sn2ff   1/1       Running   0          38s
```

To run the SATA pod:
```sh
$ kubectl create -f https://raw.githubusercontent.com/kairen/simple-device-plugin/master/artifacts/test-device-pod.yml
$ kubectl get po
NAME              READY     STATUS    RESTARTS   AGE
test-device-pod   1/1       Running   0          30s

$ kubectl exec -ti test-device-pod sh
/ # ls /dev/ | grep "sd[a-z]"
sdb
/ # mkfs.vfat /dev/sdb
/ # od -vAn -N4 -tu4 < /dev/sdb
 1838176491
```
