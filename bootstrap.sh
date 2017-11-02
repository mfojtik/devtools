#!/bin/bash

# Install base tools
brew install bash gcc jpeg s-lang bash-completion@2 socat binutils binwalk git p7zip ccat \ 
  tig coreutils curl gnu-sed lua vim diff-so-fancy midnight-commander watch \
  htop moreutils wget docker-completion xz findutilsj ripgrep

brew tap caskroom/cask
brew install keybase viscosity docker google-chrome iterm2 spotify vlc

mkdir -p ~/go/src/github.com/{mfojtik,openshift}

# standard gopath
#
cd ~/go/src/github.com/mfojtik && \
  git clone git@github.com:mfojtik/devtools.git

exit 0

cd ~/go/src/github.com/openshift && \
  git clone git@github.com:mfojtik/origin.git && \
  cd origin && \
  git remote add upstream https://github.com/openshift/origin && \
  git fetch upstream

mkdir -p ~/workspaces/kubernetes/src/k8s.io/kubernetes

cd ~/workspaces/kubernetes/src/k8s.io/kubernetes && \
  git clone git@github.com:mfojtik/kubernetes.git && \
  cd kubernetes && \
  git remote add upstream https://github.com/kubernetes/kubernetes && \
  git fetch upstream && \
  git remote add openshift https://github.com/openshift/kubernetes && \
  git fetch openshift
