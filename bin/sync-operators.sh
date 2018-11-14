#!/bin/bash

cat >${HOME}/.rsync/operator.files <<EOF
+ src
+ src/github.com
+ src/github.com/openshift
+ src/github.com/openshift/cluster-kube-apiserver-operator/***
+ src/github.com/openshift/cluster-kube-controller-manager-operator/***
+ src/github.com/openshift/cluster-kube-scheduler-operator/***
+ src/github.com/openshift/cluster-openshift-apiserver-operator/***
+ src/github.com/openshift/cluster-openshift-controller-manager-operator/***
- src/**
- bin
- pkg
EOF

rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/operator.files" ${HOME}/go/ dev:/home/mfojtik/go/