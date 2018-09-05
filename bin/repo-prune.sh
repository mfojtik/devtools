#!/bin/bash

branches=$(git branch --all --merged master | grep -v 'master$' | grep -v 'remotes/upstream' | grep -v '*')
branches_to_delete=""

set -e

for branch in ${branches[*]}; do
  name=${branch#"remotes/origin/"}
  branches_to_delete="${branches_to_delete} ${name}"
done

if [[ ! -z "${branches_to_delete}" ]]; then
  git push --quiet -d origin $branches_to_delete
fi

if [[ "$1" == "-v" ]]; then
  for branch in `git branch -r --no-merged | grep -v HEAD | grep -v 'upstream/'`; do
      echo -e `git show --format="%ci %cr %an" $branch | head -n 1` \\t$branch
  done | sort -r
fi