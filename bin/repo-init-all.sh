#!/bin/bash

set -e

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

${script_dir}/repo-add.sh openshift api
${script_dir}/repo-add.sh openshift client-go
${script_dir}/repo-add.sh openshift cluster-kube-apiserver-operator
${script_dir}/repo-add.sh openshift imagebuilder
${script_dir}/repo-add.sh openshift library-go
${script_dir}/repo-add.sh openshift origin
${script_dir}/repo-add.sh openshift openshift-docs
${script_dir}/repo-add.sh openshift release
