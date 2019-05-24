# joshuawebsite

> personal website front-end nuxt

## Build Setup

``` bash
# install dependencies
$ npm install

# serve with hot reload at localhost:3000
$ npm run dev

# build for production and launch server
$ npm run build
$ npm start

# generate static project
$ npm run generate
```

For detailed explanation on how things work, checkout [Nuxt.js docs](https://nuxtjs.org).

## Format

- autoformat vscode: <kbd>Shift</kbd> + <kbd>Option</kbd> + <kbd>F</kbd>
- lint: `npm run lint`

## Deploy

``` bash
# deploy to now.sh
$ now

# deploy to firebase
$ firebase deploy

# deploy to netlify
# merge with master branch (see below)
```

## Git Cheat sheet

``` bash
# add changes
$ git add -A

# commit changes
$ git commit -m "<message here>"

# push changes to github - frontend branch
$ git push origin frontend
```

### merge to master branch and push

``` bash
# checkout to master branch
$ git checkout master

# merge frontend branch to master
$ git merge frontend
# if there are any problems with merge they must be resolved before push
# or the push won't work

# push changes to github - master branch
$ git push origin master
```
