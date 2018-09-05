#!/bin/bash

set -e
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

dirs=(
    openshift/api
    openshift/client-go
    openshift/library-go
    openshift/cluster-kube-apiserver-operator
    openshift/origin
    openshift/release
    openshift/openshift-docs
    openshift/imagebuilder
)

for repo in ${dirs[*]}; do
    i=$((i+1)) 
    cd "${HOME}/go/src/github.com/${repo}" && ${script_dir}/repo-update.sh &
    pids[${i}]=$!
done

echo -ne "Updating ${#pids[*]} repositories "
for pid in ${pids[*]}; do
    echo -ne "."
    wait $pid
done
echo " done!"