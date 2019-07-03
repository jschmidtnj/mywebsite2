#!/usr/bin/env bash

# abort on errors
set -e

export $(grep -v '^#' .env | xargs -d '\n')

git config --global user.name $HEROKUUSERNAME
git config --global user.password $HEROKUPASSWORD

yarn build

sed -i 's/config/dummy/g' .gitignore

git init
git remote add heroku https://git.heroku.com/joshuawebsiteamp.git
git add -A
git commit -m "deploying to heroku"
git push heroku master -f

sed -i 's/dummy/config/g' .gitignore

rm -rf .git
