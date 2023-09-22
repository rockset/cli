#!/bin/sh

set -e

# TODO replace with goreleaser

go vet ./...

# check that we're on master
if [ "$(git branch --show-current)" != "master" ]; then
  echo "not on master branch"
  exit 1
fi

# check clean git repo
if [ ! -z "${git show --short --ahead-behind}" ]; then
  git show
  echo "need a clean repo"
  exit 1
fi

# extract version string
VERSION="$(grep -ohE 'v\d+\.\d+\.\d+' version.go)"
if [ -z "${VERSION}" ]; then
  echo "empty version"
  exit 1
fi

# tag repo
git tag "${VERSION}"

# push the tag
git push origin "${VERSION}"
