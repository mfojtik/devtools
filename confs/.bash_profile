# ~/.bash_profile: Settings of environment variables as well as commands that should be executed when you login

[[ -f "${HOME}/.bashrc" ]] && . "${HOME}/.bashrc"

export EDITOR="vim"     # Default editors
export VISUAL="vim"

# Custom PATH
export PATH=${HOME}/bin:${HOME}/go/src/github.com/mfojtik/devtools/bin:$PATH 

export PAGER=less
export BAT_PAGER=""
export LC_ALL=en_US.UTF-8

if [[ "$OSTYPE" == darwin* ]]; then
  export BROWSER='open'
fi

# Source bash completion
[[ -r "/usr/local/etc/profile.d/bash_completion.sh" ]] && . "/usr/local/etc/profile.d/bash_completion.sh"

if [ -f "/usr/local/opt/bash-git-prompt/share/gitprompt.sh" ]; then
  GIT_PROMPT_FETCH_REMOTE_STATUS=0
  GIT_PROMPT_IGNORE_SUBMODULES=1
  GIT_PROMPT_SHOW_UNTRACKED_FILES=no
  GIT_PROMPT_THEME=Single_line_Dark

  __GIT_PROMPT_DIR="/usr/local/opt/bash-git-prompt/share"
  source "/usr/local/opt/bash-git-prompt/share/gitprompt.sh"
fi
