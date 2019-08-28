#!/usr/local/bin/bash

set -e
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

dirs=(
    openshift/api
    openshift/library-go
    openshift/apiserver-library-go
    openshift/origin
    openshift/oc
    openshift/client-go
    openshift/openshift-apiserver
    openshift/openshift-controller-manager
    openshift/cluster-kube-apiserver-operator
    openshift/cluster-kube-controller-manager-operator
    openshift/cluster-openshift-apiserver-operator
    openshift/cluster-openshift-controller-manager-operator
    openshift/cluster-authentication-operator
    openshift/cluster-kube-scheduler-operator
    openshift/cluster-config-operator
    openshift/insights-operator
)

# require bash 4.x, hence the /usr/local/bin/bash shebang...
declare -A kube_forks
kube_forks=(
    ["kubernetes"]="openshift/kubernetes"
    ["api"]="openshift/kubernetes-api"
    ["apiextensions-apiserver"]="openshift/kubernetes-apiextensions-apiserver"
    ["apimachinery"]="openshift/kubernetes-apimachinery"
    ["apiserver"]="openshift/kubernetes-apiserver"
    ["client-go"]="openshift/kubernetes-client-go"
    ["cli-runtime"]="openshift/kubernetes-cli-runtime"
    ["cloud-provider"]="openshift/kubernetes-cloud-provider"
    ["cluster-bootstrap"]="openshift/kubernetes-cluster-bootstrap"
    ["code-generator"]="openshift/kubernetes-code-generator"
    ["component-base"]="openshift/kubernetes-component-base"
    ["csi-api"]="openshift/kubernetes-csi-api"
    ["csi-translation-lib"]="openshift/kubernetes-csi-translation-lib"
    ["kube-aggregator"]="openshift/kubernetes-kube-aggregator"
    ["kube-controller-manager"]="openshift/kubernetes-kube-controller-manager"
    ["kubelet"]="openshift/kubernetes-kubelet"
    ["kube-proxy"]="openshift/kubernetes-kube-proxy"
    ["kube-scheduler"]="openshift/kubernetes-kube-scheduler"
    ["metrics"]="openshift/kubernetes-metrics"
    ["sample-apiserver"]="openshift/kubernetes-sample-apiserver"
    ["sample-cli-plugin"]="openshift/kubernetes-sample-cli-plugin"
    ["sample-controller"]="openshift/kubernetes-sample-controller"
)

# --init cause to initialize all repos i want to keep updating later
if [[ "$1" == "--init-kube" ]]; then
  mkdir -p "${HOME}/go/src/k8s.io"
  for repo in "${!kube_forks[@]}"; do
      [[ -d ""${HOME}/go/src/k8s.io/${repo}"" ]] && continue

      pushd "${HOME}/go/src/k8s.io" >/dev/null
      git clone --origin=upstream git@github.com:kubernetes/${repo}.git
      popd >/dev/null

      pushd "${HOME}/go/src/k8s.io/${repo}" >/dev/null
      git remote add openshift "git@github.com:${kube_forks[$repo]}.git"
      git remote add mfojtik "git@github.com:mfojtik/${repo}.git"
      git fetch openshift
      popd >/dev/null
  done
  exit 0
fi


# --init cause to initialize all repos i want to keep updating later
if [[ "$1" == "--init-openshift" ]]; then
  mkdir -p "${HOME}/go/src/github.com/openshift"
  for repo in ${dirs[*]}; do
      [[ -d ""${HOME}/go/src/github.com/${repo}"" ]] && continue

      orgname=$(echo $repo | cut -d'/' -f1)
      reponame=$(echo $repo | cut -d'/' -f2)

      pushd "${HOME}/go/src/github.com/${orgname}" >/dev/null
      git clone git@github.com:mfojtik/${reponame}.git
      popd >/dev/null

      pushd "${HOME}/go/src/github.com/${repo}" >/dev/null
      git remote add upstream "git@github.com:${repo}.git"
      git fetch upstream
      popd >/dev/null
  done
  exit 0
fi

# k8s
# pushd "${HOME}/go/src/k8s.io/kubernetes" >/dev/null
# ${script_dir}/repo-update.sh
# popd >/dev/null

for repo in ${dirs[*]}; do
    pushd "${HOME}/go/src/github.com/${repo}" >/dev/null
    ${script_dir}/repo-update.sh
    popd >/dev/null
done

