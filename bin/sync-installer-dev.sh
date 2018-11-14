#!/bin/bash

rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/installer.files" ${HOME}/go/ dev:/home/mfojtik/go/
#rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/installer.files" ${HOME}/go/ 10.211.55.4:/home/mfojtik/go/
