#!/usr/local/bin/bash

if prlctl list --full | grep dev-centos >/dev/null; then
  echo "the dev is already up and running, yay!"
  exit 0
fi

prlctl start dev-centos
