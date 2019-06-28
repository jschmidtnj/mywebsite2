#!/usr/bin/env bash

# abort on errors
set -e

changes() {
  git diff --name-only --diff-filter=ADMR @~..@
}

travis_ignore=" [skip ci]"

if ! changes | grep -q "flutter" ; then
  echo "add to git commit message"
  old_commit_message=$(git log --format=%B -n1)
  echo $old_commit_message
  if echo "$old_commit_message" | grep -q "$travis_ignore" ; then
    echo "add to commit message"
    git commit --amend -m "$old_commit_message" -m "$travis_ignore"
  else
    echo "don't add to commit message"
  fi
else
  echo "dont add to git commit message"
fi
