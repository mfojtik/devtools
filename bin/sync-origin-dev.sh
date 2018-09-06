#!/bin/bash

rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/openshift.files" ${HOME}/go/ dev:/home/mfojtik/go/

# ssh dev /bin/bash -c "cd /home/mfojtik/go/src/github.com/openshift/origin && git clean -fd"
