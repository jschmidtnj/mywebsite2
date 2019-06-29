#!/usr/bin/env bash

# abort on errors
set -e

flutter pub run flutter_launcher_icons:main
flutter clean
flutter build apk
cd android && fastlane alpha

cd -
