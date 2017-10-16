#!/bin/bash

if [[ ! -z "$1" ]]; then
  docker logs origin 2>&1 | egrep --color=always -i "$1"
else
  docker logs origin --tail 1000 -f 2>&1
fi
