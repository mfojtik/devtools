#!/bin/bash

set -e

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

if [[ ! -d "$(pwd)/.git" ]]; then
  echo "$(pwd)/.git is not a git directory"
  exit 1
fi

if [[ "$(git rev-parse --abbrev-ref HEAD)" != "master" ]]; then
  echo "-> checking out the master branch"
  git checkout master &>/dev/null
fi

echo "-> fetching the latest origin remote ..."
git fetch -p origin

if [[ "$(git remote show | grep upstream || true)" != "upstream" ]]; then
  echo "$(pwd) does not have upstream remote defined"
  exit 0
fi

echo "-> fetching the latest upstream remote ..."
git fetch -p upstream

echo "-> reseting local master branch to upstream/master"
git reset --hard upstream/master

echo "-> pushing local master branch to origin remote ..."
git push origin master --force

echo "-> deleting local branches already merged to master ..."
# Prune branches already merged to upstream master
git branch --merged master | grep -v 'master$' | xargs git branch -D &>/dev/null
