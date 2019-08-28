#!/usr/local/bin/bash

bold=$(tput bold)
normal=$(tput sgr0)

function report() {
  if [[ -f /usr/local/bin/emojify ]]; then
    emojify $@
  else
    echo $@
  fi
}

function latest_commit() {
  echo -n $(git log --pretty=oneline --no-decorate --abbrev-commit --no-merges -n1 origin/master)
}

function new_commits() {
  local current="${1}"
  local latest="${2}"
  git log --pretty=format:"   %h: %ar: %s" --no-decorate --abbrev-commit --no-merges ${current}..${latest}
}

function current_branch_name() {
  echo -n $(git rev-parse --abbrev-ref HEAD)
}

function origin_remote_url() {
  echo -n $(git remote get-url origin)
}

set -e

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

report ":bricks: Processing repository ${bold}$(origin_remote_url)${normal} ..."

if [[ ! -d "$(pwd)/.git" ]]; then
  report ":poop:  $(pwd)/.git is not a git directory"
  exit 1
fi

report ":cloud: Fetching remote updates for: $(echo -n $(git remote)) ..."
git fetch --quiet --all --prune --no-tags

if [[ "$(git remote show | grep upstream || true)" != "upstream" ]]; then
  report ":poop:  $(pwd) does not have 'upstream' remote" && exit 0
fi

current_master_rev=$(git rev-parse master)
current_local_rev=$(git rev-parse origin/master)
current_upstream_rev=$(git rev-parse upstream/master)

if [[ "${current_local_rev}" != "${current_master_rev}" ]]; then
  report ":hammer_and_wrench: Found local master and origin/master diverge, fixing local"
  git reset --quiet --hard origin/master
fi

if [[ "${current_upstream_rev}" == "${current_local_rev}" ]]; then
  report ":sunny: Everything is up-to-date for ${bold}$(origin_remote_url)${normal}: $(latest_commit)" && echo
  exit 0
fi

if [[ "$(current_branch_name)" != "master" ]]; then
  report ":cloud_tornado: Current branch is ${bold}'$(current_branch_name)'${normal}, switching to 'master'"
  set +e; git checkout --quiet master; set -e
  [[ "$(current_branch_name)" != "master" ]] && report ":poop: Unable to switch $(basename $(pwd)) to master branch." && echo && exit 0
fi

report ":hammer_and_wrench: Updating origin/master to point to upstream/master"
git reset --quiet --hard upstream/master

report ":cloud: Updating remote origin/master branch ..."
git push --quiet origin master --force

report ":recycle: Cleaning up local branches merged in upstream ..."
git branch --merged master | grep -v 'master$' | xargs git branch -D

report ":sunny: ${bold}$(origin_remote_url)${normal} is now $(latest_commit)"
new_commits ${current_local_rev} ${current_upstream_rev} && echo
