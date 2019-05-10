#!/bin/sh
set -ex
docker build -t gomkbuild .
docker run -it -v`pwd`:/gomkbuild -w/gomkbuild gomkbuild                       \
  go test -v -coverprofile=gomkbuild.cov ./...
