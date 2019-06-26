#!/usr/bin/env bash

# abort on errors
set -e

flutter clean
flutter build apk
cd android && fastlane alpha

cd -
