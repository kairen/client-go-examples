#!/bin/bash

# Copyright © 2018 Kyle Bai(kyle.b@inwinstack.com)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

# Add golang repository
sudo add-apt-repository -y ppa:hnakamur/golang-1.9

# Install docker engine
curl -fsSL "https://get.docker.com/" | sh

# Install kubernetes
curl -s "https://packages.cloud.google.com/apt/doc/apt-key.gpg" | sudo apt-key add -
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update && sudo apt-get install -y kubelet kubeadm kubectl kubernetes-cni golang-go golang-1.9-doc

# Set swap
sudo swapoff -a && sudo sysctl -w vm.swappiness=0

# Create src
PROJECT_HOME="/root/go/src/github.com/kairen/simple-operator"
sudo mkdir -p ${PROJECT_HOME}

echo 'export GOPATH=$HOME/go' | sudo tee -a /root/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' | sudo tee -a /root/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /root/.bashrc

sudo mkdir -p /etc/cni/net.d /root/.kairen/certs /root/.kairen/certs/ /root/.kube/
sudo cp ${PROJECT_HOME}/hack/k8s/10-kubeadm.conf /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
sudo cp ${PROJECT_HOME}/hack/k8s/config.yml /var/lib/config.yml
sudo cp ${PROJECT_HOME}/hack/k8s/10-net.conf /etc/cni/net.d/10-net.conf
sudo systemctl daemon-reload
sudo kubeadm init --config /var/lib/config.yml --skip-preflight-checks
sudo cp /etc/kubernetes/admin.conf /root/.kube/config

sudo kubectl apply -f ${PROJECT_HOME}/artifacts/operator/operator-rbac.yml

SECRET_NAME=$(sudo kubectl get sa operator -o yaml | grep -o "operator-token-\w*")
TOKEN=$(sudo kubectl get secret ${SECRET_NAME} -o go-template='{{.data.token}}')
echo "export TOKEN=${TOKEN}" | sudo tee -a /root/tokenrc
