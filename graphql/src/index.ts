import express = require('express')
import bodyParser = require('body-parser')
import cors = require('cors')
import cookieParser = require('cookie-parser')
import expressSession = require('express-session')
import { ApolloServer, gql, AuthenticationError } from 'apollo-server-express'
import passport from 'passport'
import jwt from 'jsonwebtoken'
import passportCustom from 'passport-custom'
import config from './config'
import { OAuth2Client } from 'google-auth-library'

const googleoauth = new OAuth2Client(config.auth.google.client_id)
const CustomStrategy = passportCustom.Strategy

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    hello: String
  }
`

// Provide resolver functions for your schema fields
const resolvers = {
  Query: {
    hello: () => 'Hello world!',
  }
}

const server = new ApolloServer({
  typeDefs,
  resolvers,
  tracing: true
})

const context = async ({ req }) => {
  const strategy = req.headers.strategy
  if (strategy) {
    switch (strategy) {
      case 'google':
        return await googleoauth.getTokenInfo(req.headers.authorization.split(' ')[1])
          .then(user => {
            return { user }
          }).catch(err => {
            throw new AuthenticationError('error getting token data from google')
            return
          })
        break
      default:
        throw new AuthenticationError('invalid authentication strategy')
        break
    }
  }
}

const app = express()

app.use(cookieParser())
app.use(bodyParser.urlencoded({
  extended: true
}))
app.options('*', cors({
  origin: true
}))
app.use(expressSession({
  secret: 'test123',
  resave: true,
  saveUninitialized: true
}))

app.use(server.graphqlPath, (req, res, next) => {
  const strategy = req.headers.strategy
  if(strategy) {
    // @ts-ignore
    if (validstrategies.includes(strategy)) {
      passport.authenticate(strategy, (err, user) => {
        if (err) {
          console.log(`passport error ${err}`)
          res.status(401).send(`token error ${err}`)
        } else {
          if (user) {
            console.log(JSON.stringify(user))
            // @ts-ignore
            req.user = user
            next()
          } else {
            res.status(401).send('could not get user data')
          }
        }
      })(req, res, next)
    } else {
      res.status(401).send('invalid strategy header found')
    }
  } else {
    next()
  }
})

app.use(passport.initialize())
app.use(passport.session())

server.applyMiddleware({ app })

const PORT = process.env.PORT || 8080

app.listen({ port: PORT }, () =>
  console.log(`Server listening on port ${PORT}, path ${server.graphqlPath} 🚀`)
)
