#!/bin/bash

# abort on errors
set -e

yarn install
cd flutter
flutter pub get
flutter pub run flutter_launcher_icons:main
cd ../nuxt
yarn install
cd ../docs
yarn install
cd ../amp
yarn install 
cd ../electron
yarn install
cd ../init
yarn install
cd ../shortlink
yarn install
cd ..
