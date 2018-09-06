#!/bin/bash

mkdir -p /Users/mfojtik/go/src/github.com/openshift/origin/_output/local/bin/linux/amd64
scp -r dev:go/src/github.com/openshift/origin/_output/local/bin/linux/amd64/* \
       /Users/mfojtik/go/src/github.com/openshift/origin/_output/local/bin/linux/amd64
