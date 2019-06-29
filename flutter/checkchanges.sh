#!/usr/bin/env bash

# abort on errors
set -e

changes() {
  git diff --name-only --diff-filter=ADMR @~..@
}

travis_ignore="[skip ci]"

if ! changes | grep "flutter/*" ; then
  old_commit_message=$(git log --format=%B -n1)
  echo $old_commit_message
  if echo "$old_commit_message" | grep -q "$travis_ignore" ; then
    echo "add skip ci to commit message"
    git commit --amend -m "$old_commit_message" -m "$travis_ignore"
  else
    echo "skip ci already added"
  fi
else
  echo "no flutter changes found"
fi
