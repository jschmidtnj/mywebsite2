#!/bin/bash

# abort on errors
set -e

yarn build

sed -i -e 's/config/dummy/g' .gitignore
rm -rf .gitignore-e

git init
git remote add heroku https://git.heroku.com/joshuawebsiteamp.git
git add -A
git commit -m "deploying to heroku"
git push heroku master -f

sed -i -e 's/dummy/config/g' .gitignore
rm -rf .gitignore-e

rm -rf .git
