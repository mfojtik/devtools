#!/bin/bash

set -e
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

if [[ ! -d "$(pwd)/.git" ]]; then
  echo "$(pwd)/.git is not a git directory"
  exit 1
fi

if [[ "$(git rev-parse --abbrev-ref HEAD)" == "master" ]]; then
  git checkout master &>/dev/null
fi

git fetch -p --quiet origin

if [[ "$(git remote show | grep upstream || true)" != "upstream" ]]; then
  echo "$(pwd) does not have upstream remote defined"
  exit 0
fi

git fetch --quiet -p upstream
git reset --hard upstream/master &>/dev/null
git push --quiet origin master --force

# Prune branches already merged to upstream master
git branch --merged master | grep -v 'master$' | xargs git branch -D &>/dev/null

# Prune branches in remote
${script_dir}/repo-prune.sh