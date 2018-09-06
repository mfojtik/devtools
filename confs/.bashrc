# .bashrc

# Source global definitions
if [ -f /etc/bashrc ]; then
  . /etc/bashrc
fi

shopt -s cdspell        # This will correct minor spelling errors in a cd command.
shopt -s histappend     # Append to history rather than overwrite
shopt -s dotglob        # Files beginning with . to be returned in the results of path-name expansion.
shopt -s nocaseglob
shopt -s checkwinsize
shopt -s histappend

HISTSIZE=500000
HISTFILESIZE=100000
HISTCONTROL="erasedups:ignoreboth"
HISTIGNORE="&:[ ]*:exit:ls:bg:fg:history:clear"
HISTTIMEFORMAT='%F %T '

export EDITOR="vim"
export VISUAL="vim"
export PATH=$HOME/bin:$PATH
export PS1="\w \[\033[00m\]\[\033[00;34m\]→\[\033[00m\] "

[ -f ~/.fzf.bash ] && source ~/.fzf.bash