#!/bin/sh
set -ex
gorepopath=/go/src/github.com/measurement-kit/go-measurement-kit
docker run -it -v `pwd`:$gorepopath openobservatory/mk-alpine:20190509         \
  $gorepopath/script/build-alpine.sh
