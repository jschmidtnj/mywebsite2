#!/usr/bin/env bash

# abort on errors
set -e

changes() {
  git diff --stat --cached -- flutter/ nuxt/
}

travis_ignore="[skip ci]"

if ! changes | grep -E "flutter/|nuxt/" ; then
  echo "no flutter changes found"
  sed -i.bak -e "1s/^/$travis_ignore /" "$GIT_DIR/COMMIT_EDITMSG"
else
  echo "flutter changes found"
fi
