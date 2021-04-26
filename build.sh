#!/bin/bash

function check_for_gofmt {
  bad=$(find cmd pkg dom -iname "*.go" | xargs gofmt -l)
  if [ -n "$bad" ]; then
    echo "gmfmt returned style errors"
    for line in $bad; do
      echo "style error found in: $line"
    done
    exit 1
  fi

  return 0
}

go version

check_for_gofmt

set -e
set -o pipefail

cmd="$1"
[ -z $1 ] && cmd='...'

if [ "${cmd}" == '...' ] || [ -d "./cmd/$cmd" ]; then
  go install -ldflags="-s -w" ./cmd/$cmd
fi

