#!/bin/bash

images=(
	cluster-kube-apiserver-operator
	cluster-kube-controller-manager-operator
	cluster-openshift-apiserver-operator
	cluster-openshift-controller-manager-operator
)

set -e
echo "-> Fetching image details from api.ci.openshift.org ..."
echo

for image in "${images[@]}"; do
	rm -f "/tmp/${image}.json"
	oc get istag -n openshift origin-v4.0:${image} --server https://api.ci.openshift.org -o json >/tmp/${image}.json
	commit=$(cat /tmp/${image}.json | jq -r '.image.dockerImageMetadata.Config.Labels."io.openshift.build.commit.id"')
	curl -s "https://api.github.com/repos/openshift/${image}/commits/${commit}?format=json" -o /tmp/${image}-commit.json
done

for image in "${images[@]}"; do
	commit=$(cat /tmp/${image}.json | jq -r '.image.dockerImageMetadata.Config.Labels."io.openshift.build.commit.id"')
	commit_message=$(cat /tmp/${image}-commit.json | jq -r '.commit.message' | xargs echo -n)
	commit_message_trim=${commit_message#"Merge pull request "}
	author=$(cat /tmp/${image}-commit.json | jq -r '.commiter.login')
	printf "%-50s %-8s %-90s\n" "${image}" "${commit::8}" "PR: ${commit_message_trim::50}..."
done
