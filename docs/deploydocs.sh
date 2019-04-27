#!/usr/bin/env bash

# abort on errors
set -e

# build
npm run build

rm -rf .vuepress/dist/

mv mywebsite2/ .vuepress/dist

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