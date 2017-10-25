#!/bin/bash

cd ${HOME}/go/src/github.com/openshift/origin

if [[ "$1" == "--force" ]]; then
  export OS_BUILD_ENV_VOLUME_FORCE_NEW=TRUE hack/env
fi

export OS_BUILD_ENV_PRESERVE="_output"
export OS_BUILD_ENV_REUSE_VOLUME="local"
export OS_BUILD_ENV_DOCKER_ARGS='-e OS_VERSION_FILE= '

current_branch=$(git rev-parse --abbrev-ref HEAD)

echo "+ Building '${current_branch}' ..."
exec ./hack/env ./hack/build-go.sh cmd/openshift
