#!/usr/bin/env bash

# abort on errors
set -e

if ! git diff --name-only $TRAVIS_COMMIT_RANGE | grep nuxt/
then
  echo "No nuxt files updated."
  return 0
fi
