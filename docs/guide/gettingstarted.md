# Getting Started

## clone

First clone the github repository: `git clone https://github.com/jschmidtnj/mywebsite2`. There are many files in a few different folders - let's break that down. The `nuxt/` folder contains the actual website source code, written in Vue.js and Typescript using Nuxt as the overarching framework. The `docs/` folder contains the documentation you are reading now. The `amp/` folder contains a Node.js server deployed to Heroku currently, and stands for Accelerated Mobile Pages (fast page cacheing from Google). The amp pages used currently are blog and project posts. The `electron/` folder contains a script for converting the generated website from nuxt to a desktop app, on mac, linux, and windows. The `flutter/` folder contains the source code for the mobile app, written using the flutter framework. Finally, all of these endpoints utilize the api created in the `graphql/` folder, which is currently deployed on google app engine. This api was written in golang and is designed to be fast and stateless, with a mongodb database and authentication using json web tokens (jwt). If you understand what all of that meant, then you're honestly good to go in developing with this stack. If not, just read a bit on each of the technologies and get an idea of what's going on.

## configuration

You need to make a few config files in order to get started. The first one you should make is the `.env` file in the `nuxt/` directory. Copy the `.default.env` file and change the parameters. You should change the APIURL parameter immediately to `localhost:8000` or whatever port you are using to develop the graphql api, and then later change it to the production url. Then change the AMPURL when you are working on getting amp going. The AUTHCONFIG parameter comes from google oath2, and can be changed later. And the seoconfig can also be updated later. Run `yarn dev` in the `nuxt/` directory to start the development server for the website, and you should be good to go.  

In the `graphql/` directory, you also need a `.env` file. Start by copying the default. The SECRET is the jwt secret - this should be changed. The port is the port the api will be running on in development. The MONGOURI is used to connect to mongodb - I used MongoDB Atlas for managed storage, but you can use whatever you want. The token expiration is in hours and is how long each jwt lasts. The sendgrid api key is used for authenticating sending emails with sendgrid - reset password, etc., and should be changed. Note that the templates for these emails come from the deployed Nuxt website, so it should be deployed for that endpoint to work. The websiteurl is the deployed url of the nuxt website. The elasticuri is the uri for elasticsearch - used for searching for blogs and projects - I used bonsai.io for managed hosting. The storageconfig is the serviceaccount json from google cloud storage for storing files, and the storagebucketname is the bucket containing all the pictures.

In the `flutter/` directory, you need to make a `config.dart` file following the sample, specifying the api url, which you can use as localhost if you are running locally to test.

In the `amp/` directory, you need to follow the `sampleconfig.ts` file, specifying the api url and website url. Same with the `init/` directory.

## test + debug

- run `yarn dev` for testing the website, in the `nuxt/` directory
- run `go run *.go` for running the api, in the `graphql/` directory
- run `flutter run` in the `flutter/` directory, for testing the mobile app
- run `yarn start` in the `amp/` directory, for testing amp
- run `yarn start` in the `init/` directory, for adding admin users to the api and configuring the api
