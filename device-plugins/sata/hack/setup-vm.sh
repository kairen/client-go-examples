#!/bin/bash
#
# Setup vagrant vms.
#

set -eu

# Copy hosts info
cat <<EOF > /etc/hosts
127.0.0.1	localhost
127.0.1.1	vagrant.vm	vagrant

192.16.35.11 k8s-m1
192.16.35.12 k8s-n1
192.16.35.13 k8s-n2

# The following lines are desirable for IPv6 capable hosts
::1     localhost ip6-localhost ip6-loopback
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
EOF

# Install docker
curl -fsSL "https://get.docker.com/" | sh

# Install kubernetes
if [ ${HOSTNAME} != "ldap-server" ]; then
  curl -s "https://packages.cloud.google.com/apt/doc/apt-key.gpg" | sudo apt-key add -
  echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
  sudo apt-get update && sudo apt-get install -y kubelet kubeadm kubectl kubernetes-cni
fi

swapoff -a && sysctl -w vm.swappiness=0
sed '/vagrant--vg-swap_1/d' -i  /etc/fstab

if [ ${HOSTNAME} == "k8s-m1" ]; then
  kubeadm init --service-cidr 10.96.0.0/12 \
    --kubernetes-version v1.10.0 \
    --pod-network-cidr 10.244.0.0/16 \
    --token b0f7b8.8d1767876297d85c \
    --apiserver-advertise-address 192.16.35.11

  # copy k8s config
  mkdir ~/.kube && cp /etc/kubernetes/admin.conf ~/.kube/config

  # deploy calico network
  kubectl apply -f "https://kairen.github.io/files/k8s-ldap/calico.yml.conf"
elif [[ ${HOSTNAME} =~ k8s-n ]]; then
  kubeadm join 192.16.35.11:6443 \
    --token b0f7b8.8d1767876297d85c \
    --discovery-token-unsafe-skip-ca-verification
fi
