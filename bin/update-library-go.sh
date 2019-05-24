#!/usr/bin/env bash

set -x
set -e

GOPATH=${HOME}/go

repos=(
    cluster-openshift-apiserver-operator
    cluster-kube-apiserver-operator
    cluster-kube-controller-manager-operator
)

rm -f /tmp/bump-library-go-urls

for repo in "${repos[@]}"; do
    echo "Updating ${repo} ..."
    pushd ${GOPATH}/src/github.com/openshift/$repo

    # Update upstream
    git fetch upstream

    # Delete bump-library-go branch if exists locally
    git checkout --force master
    git branch -D bump-library-go || true

    # Checkout
    git checkout --force upstream/master -b bump-library-go
    make update-deps

    # Add all modified files and commit
    git add -A
    git commit -m "bump(*): library-go"
    git push -f origin bump-library-go

    # Create pull request
    hub pull-request -m "bump library-go" >>/tmp/bump-library-go-urls
    popd
done

echo
echo
echo "Pull request URLs:"
cat /tmp/bump-library-go-urls
