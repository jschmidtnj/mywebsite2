#!/bin/bash

# abort on errors
set -e

npm run precommit --prefix nuxt
gofmt -w graphql
flutter format flutter
npm run precommit --prefix amp
npm run precommit --prefix electron
npm run precommit --prefix init
npm run precommit --prefix docs
npm run precommit --prefix shortlink
git add -A
