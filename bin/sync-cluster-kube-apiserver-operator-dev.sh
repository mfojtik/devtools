#!/bin/bash

rsync --info=progress2  -az --delete --filter=". ${HOME}/.rsync/cluster-kube-apiserver-operator.files" ${HOME}/go/ dev:/home/mfojtik/go/
