#!/bin/bash
set -e

files=(
    .bash_profile
    .bashrc
    .gitconfig
    .gitignore
    .go-vars
    .go-vars.default
    .go-vars.kube
    .go-vars.ose
    .hushlogin
    .inputrc
    .tigrc
)

pushd ${HOME}
for f in ${files[*]}; do
  rm ${HOME}/go/src/github.com/mfojtik/devtools/confs/${f}
  cp -f -v ${f} ${HOME}/go/src/github.com/mfojtik/devtools/confs/
done
popd