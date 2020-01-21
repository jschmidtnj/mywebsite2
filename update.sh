#!/bin/bash

# abort on errors
set -e

yarn upgrade
cd amp
yarn upgrade
cd ../docs
yarn upgrade
cd ../electron
yarn upgrade
cd ../init
yarn upgrade
cd ../nuxt
yarn upgrade
cd ../shortlink
yarn upgrade
cd ../graphql
go get -u
cd ../flutter/android
bundle update
cd ../ios
bundle update
cd ../..
