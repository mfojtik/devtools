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
)

for repo in ${dirs[*]}; do
    echo "[*] Updating ${repo} ..."
    pushd "${HOME}/go/src/github.com/${repo}" >/dev/null
    ${script_dir}/repo-update.sh
    popd >/dev/null
done

echo
echo "All repositories successfully updated!"
