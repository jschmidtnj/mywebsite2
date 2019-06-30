#!/usr/bin/env bash

# abort on errors
set -e

changes() {
  git diff --stat --cached -- flutter/
}

travis_ignore="[skip ci]"

if ! changes | grep "flutter/*" ; then
  echo "no flutter changes found"
  BRANCH_NAME=$(git branch | grep '*' | sed 's/* //')
  sed -i.bak -e "1s/^/$travis_ignore /" "$GIT_DIR/COMMIT_EDITMSG"
else
  echo "flutter changes found"
fi
