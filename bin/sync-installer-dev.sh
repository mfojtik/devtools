#!/bin/bash

rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/installer.files" ${HOME}/go/ dev:/home/mfojtik/go/
