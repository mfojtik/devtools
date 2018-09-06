#!/bin/bash

set -e

# Install base tools
brew install bash gcc jpeg s-lang bash-completion@2 socat binutils binwalk git p7zip ccat \ 
  tig coreutils curl gnu-sed lua vim diff-so-fancy midnight-commander watch \
  htop moreutils wget docker-completion xz findutilsj ripgrep

brew tap caskroom/cask
brew install keybase viscosity iterm2 vlc

mkdir -p ~/go/src/github.com/{mfojtik,openshift}

# standard gopath
#
cd ${HOME}/go/src/github.com/mfojtik && \
  git clone git@github.com:mfojtik/devtools.git

pushd ${HOME}/go/src/github.com/mfojtik/devtools
cp -v confs/.* ${HOME}/
popd