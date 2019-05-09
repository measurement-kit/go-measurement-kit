#!/bin/sh
set -ex
dirname=`dirname $0`
topdir=`cd $dirname && pwd -P`
cd $topdir/..
export GOPATH=/go
apk add --no-progress git go
go get -v ./...
go test -v ./...
