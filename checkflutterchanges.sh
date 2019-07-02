#!/usr/bin/env bash

# abort on errors
set -e

if ! git diff --name-only $TRAVIS_COMMIT_RANGE | grep flutter/ | grep -vE '(.md)'
then
  echo "No flutter files updated."
  return 0
fi
