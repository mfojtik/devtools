#!/bin/bash

set -e
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

dirs=(
    openshift/api
    openshift/library-go
    openshift/origin
    openshift/installer
    openshift/cluster-kube-apiserver-operator
    openshift/cluster-kube-controller-manager-operator
    openshift/cluster-openshift-apiserver-operator
    openshift/cluster-openshift-controller-manager-operator
    openshift/cluster-kube-scheduler-operator
    openshift/cluster-config-operator
)

# kube repo
pushd "${HOME}/go/src/k8s.io/kubernetes" >/dev/null
${script_dir}/repo-update.sh
popd >/dev/null

for repo in ${dirs[*]}; do
    pushd "${HOME}/go/src/github.com/${repo}" >/dev/null
    ${script_dir}/repo-update.sh
    popd >/dev/null
done

