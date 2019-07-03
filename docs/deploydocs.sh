#!/usr/bin/env bash

# abort on errors
set -e

export $(grep -v '^#' .env | xargs -d '\n')

git config --global user.name $GITHUBUSERNAME
git config --global user.password $GITHUBPASSWORD

# build
yarn build
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

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

git init
git add -A
git commit -m 'deploy'

# if you are deploying to https://<USERNAME>.github.io
# git push -f git@github.com:<USERNAME>/<USERNAME>.github.io.git master

# if you are deploying to https://<USERNAME>.github.io/<REPO>
# git push -f git@github.com:<USERNAME>/<REPO>.git master:gh-pages
git push -f https://github.com/jschmidtnj/mywebsite2.git master:gh-pages

cd -
