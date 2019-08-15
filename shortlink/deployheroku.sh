#!/bin/bash

# abort on errors
set -e

yarn precommit

sed -i -e 's/.env/dummy/g' .gitignore
rm -rf .gitignore-e

git init
git remote add heroku https://git.heroku.com/joshuashortwebsite.git
git add -A
git commit -m "deploying to heroku"
git push heroku master -f

sed -i -e 's/dummy/.env/g' .gitignore
rm -rf .gitignore-e

rm -rf .git
