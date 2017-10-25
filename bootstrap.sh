#!/bin/bash

# standard gopath
mkdir -p ~/go/src/github.com/{mfojtik,openshift}

cd ~/go/src/github.com/mfojtik && \
  git clone git@github.com:mfojtik/devtools.git

cd ~/go/src/github.com/openshift && \
  git clone git@github.com:mfojtik/origin.git && \
  cd origin && \
  git remote add upstream https://github.com/openshift/origin && \
  git fetch upstream

mkdir -p ~/workspaces/kubernetes/src/k8s.io/kubernetes

cd ~/workspaces/kubernetes/src/k8s.io/kubernetes && \
  git clone git@github.com:mfojtik/kubernetes.git && \
  cd kubernetes && \
  git remote add upstream https://github.com/kubernetes/kubernetes && \
  git fetch upstream && \
  git remote add openshift https://github.com/openshift/kubernetes && \
  git fetch openshift
