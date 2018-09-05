# Custom bash configuration
#
#
source ${HOME}/.bashrc
[[ -f "${HOME}/.tokens" ]] && source ${HOME}/.tokens

# Enable bash completion
if [ -f /usr/local/share/bash-completion/bash_completion ]; then
  . /usr/local/share/bash-completion/bash_completion
fi

bind "set completion-ignore-case on"
bind "set completion-map-case on"

# Enable gvm
[[ -s "${HOME}/.gvm/scripts/gvm" ]] && source "${HOME}/.gvm/scripts/gvm"

# Unlock SSH key
ssh-add -K ~/.ssh/id_rsa &>/dev/null

# Detect which `ls` flavor is in use
if ls --color > /dev/null 2>&1; then # GNU `ls`
  colorflag="--color"
else # OS X `ls`
  colorflag="-G"
fi

export PATH="/usr/local/bin:$PATH"
export LANG='en_US.UTF-8'
export LC_ALL='en_US.UTF-8'
export HOMEBREW_NO_ANALYTICS=1
export PROMPT_DIRTRIM=2

export GIT_AUTHOR_NAME='Michal Fojtik'
export GIT_COMMITTER_NAME='Michal Fojtik'
export GIT_AUTHOR_EMAIL='mfojtik@redhat.com'
export GIT_COMMITTER_EMAIL='mfojtik@redhat.com'
export CLICOLOR=1
export EDITOR='vim'

# This allow to test origin using unsupported go version
export PERMISSIVE_GO=y

# Want to see all processes and cores
alias htop='sudo htop'
alias dmesg='sudo dmesg'
alias cat='ccat'

# Lazy switching my brain
alias ack='rg'

# Always use color output for `ls`
alias ls="command ls ${colorflag}"

# Force push to update github pulls
alias gitpush="git push -f"

# Include my env scripts
export PATH="${HOME}/go/src/github.com/mfojtik/devtools/env:$PATH"

os-get-memory-profile() {
  oc exec -c prom-proxy prometheus-0 -- /bin/bash -c 'curl -k https://52.60.163.81:8444/debug/pprof/profile -H "Authorization: bearer $( cat /var/run/secrets/kubernetes.io/serviceaccount/token )"' > /tmp/profile
}

# Allow easy switching between different golang workspaces
switch-go-workspace() {
  local workspace="$1"
  local defaultdir="$2"
  ln -sf "${HOME}/.go-vars.${workspace}" "${HOME}/.go-vars"
  source "${HOME}/.go-vars"
  if [ ! -d "$GOPATH" ]; then
    echo "--> ERROR: $GOPATH does not exists!"
    return 1
  fi
  echo -e "[\e[1m\033[32m${workspace}\e[0m] GOPATH=\033[92m${GOPATH}\e[0m" && cd "${GOPATH}/src/${defaultdir}"
}

cd-origin() {
  switch-go-workspace "default" "github.com/openshift/origin"
}

cd-mfojtik() {
  switch-go-workspace "default" "github.com/mfojtik"
}

cd-kube() {
  switch-go-workspace "kube" "k8s.io/kubernetes"
}

cd-ose() {
  switch-go-workspace "ose" "github.com/openshift/ose"
}

test -e ${HOME}/.go-vars && source ${HOME}/.go-vars

alias run-dev="${HOME}/vms/centos7/run_dev.sh"
alias kill-dev="sudo pkill -9 xhyve"

[[ -s "/Users/mfojtik/.gvm/scripts/gvm" ]] && source "/Users/mfojtik/.gvm/scripts/gvm"
export PATH="/usr/local/opt/curl/bin:$PATH"

if which hub &>/dev/null; then
  eval "$(hub alias -s)"
fi
