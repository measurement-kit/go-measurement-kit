#!/bin/sh
set -ex
brew tap measurement-kit/measurement-kit
brew update
brew upgrade
brew install measurement-kit
go test -v -coverprofile=gomkbuild.cov ./...
