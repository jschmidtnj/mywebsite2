# joshuawebsite

> personal website front-end nuxt

## Build Setup

``` bash
# install dependencies
$ npm install

# serve with hot reload at localhost:3000
$ yarn dev

# build for production and launch server
$ yarn build
$ npm start

# generate static project
$ yarn generate
```

For detailed explanation on how things work, checkout [Nuxt.js docs](https://nuxtjs.org).

## Format

- autoformat vscode: <kbd>Shift</kbd> + <kbd>Option</kbd> + <kbd>F</kbd> (windows)
- for linux it's <kbd>Ctrl</kbd> + <kbd>Shift</kbd> + <kbd>I</kbd>
- lint: `yarn lint`

## Deploy

``` bash
# deploy to now.sh
$ now

# deploy to firebase
$ firebase deploy

# deploy to netlify
# merge with master branch (see below)
```

## now secrets

`now secret add key value`

add secret without escaped characters in json

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

## now.sh and cloudflare

see [this](https://zeit.co/docs/v1/guides/how-to-use-cloudflare) for more info.  
now.sh doesn't work currently because we need an ssr config. instead I'm using heroku (or you can use gcp)

## heroku

deploy to heroku using the script. see [this](https://support.cloudflare.com/hc/en-us/articles/205893698-Configure-Cloudflare-and-Heroku-over-HTTPS) for connecting to cloudflare. also look at [this](https://nuxtjs.org/faq/heroku-deployment/) first.
