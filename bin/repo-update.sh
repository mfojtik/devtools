#!/bin/bash

set -e

if [[ ! -d "$(pwd)/.git" ]]; then
  echo "$(pwd)/.git is not a git directory"
  exit 1
fi

if [[ "$(git rev-parse --abbrev-ref HEAD)" == "master" ]]; then
  git checkout master &>/dev/null
fi

git fetch --quiet origin
git fetch --quiet upstream
git reset --hard upstream/master &>/dev/null
git push --quiet origin master --force