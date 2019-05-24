# ~/.bashrc executed by bash(1) for non-login shells.

# If not running interactively, don't do anything
case $- in
    *i*) ;;
      *) return;;
esac

# Source OSX defaults
if [ -f /etc/bashrc ]; then
  . /etc/bashrc
fi

shopt -s cdspell        # This will correct minor spelling errors in a cd command.
shopt -s histappend     # Append to history rather than overwrite
shopt -s dotglob        # Files beginning with . to be returned in the results of path-name expansion.
shopt -s globstar       # The ** pattern will match all files and zero or more directories

HISTSIZE=1000           # Maximum lines in history file
HISTFILESIZE=2000       # Maximum history file size
HISTCONTROL=ignoreboth  # Do not put duplicate lines to history

# Generic
alias ll='ls -latr'

# Quick path switch
alias cd-origin="cd ${HOME}/go/src/github.com/openshift/origin"
alias cd-mfojtik="cd ${HOME}/go/src/github.com/mfojtik"
alias cd-library-go="cd ${HOME}/go/src/github.com/openshift/library-go"

# ¯\_(ツ)_/¯
alias ack="rg"

if [[ -f "/usr/local/bin/bat" ]]; then
    alias cat="bat -p"
fi

if [[ "$OSTYPE" == darwin* ]]; then
    alias sed="gsed"
    alias aws="gawk"
fi

# vs-code
function code {
    if [[ $# = 0 ]]
    then
        open -a "Visual Studio Code"
    else
        local argPath="$1"
        [[ $1 = /* ]] && argPath="$1" || argPath="$PWD/${1#./}"
        open -a "Visual Studio Code" "$argPath"
    fi
}