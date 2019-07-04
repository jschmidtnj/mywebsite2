#!/usr/bin/env bash

# abort on errors
set -e

# build
yarn generate
yarn pwa
sed -i 's/\/images\//\/mywebsite2\/images\//g' docs/Polyfills/manifest.json

rm -rf .vuepress/dist/
mv mywebsite2/ .vuepress/dist
rm -rf docs/Polyfills/web/Images/mywebsite2
cp -r docs/Polyfills/web/Images/ .vuepress/dist/images
rm -rf docs/Polyfills/web/Images/
mv docs/Polyfills/manifest.json .vuepress/dist/manifest.json
rm -rf docs/

# navigate into the build output directory
cd .vuepress/dist

# make sure to create secret file with following commands:
# tar cvf secrets.tar .env
# travis encrypt-file secrets.tar
