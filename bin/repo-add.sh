#!/bin/bash

set -e

orgname="${1}"
reponame="${2}"

# Change this if the username != github username
forkname="$(whoami)"

dstdir="${HOME}/go/src/github.com/${orgname}"

if [[ -d "${dstdir}/${reponame}" ]]; then
    echo "directory ${dstdir}/${reponame} already exists"
    exit 1
fi

mkdir -p ${dstdir}
pushd "${dstdir}"
git clone git@github.com:${forkname}/${reponame}.git
popd

if [[ "${orgname}" != "${forkname}" ]]; then
  pushd "${dstdir}/${reponame}"
  git remote add upstream https://github.com/${orgname}/${reponame}
  git fetch upstream
  popd
fi