#!/bin/bash

cd $HOME/go/src/github.com/openshift/origin

current_version=$(git describe --abbrev=0 --tags)
data_dir="${HOME}/.config/openshift/${current_version}"
binary_path="_output/local/bin/linux/amd64/openshift"


if [[ ! -x "$(pwd)/${binary_path}" ]]; then
  echo "ERROR: OpenShift binary not found in $(pwd)/${binary_path}"
  exit 1
fi

mkdir -p ${data_dir}

export OPENSHIFT_MASTER_EXTENDED_ARGS="--loglevel=4"

exec oc-dev cluster up --version="${current_version}" \
  --openshift-binary-path=$(pwd)/${binary_path} \
  --host-data-dir="$data_dir" \
  --host-config-dir="${data_dir}/config" \
  --use-existing-config \
  --image-streams=centos7 \
  --skip-registry-check
