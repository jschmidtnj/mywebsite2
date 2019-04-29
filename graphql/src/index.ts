import express = require('express')
import bodyParser = require('body-parser')
import cors = require('cors')
import cookieParser = require('cookie-parser')
import { ApolloServer, gql, AuthenticationError } from 'apollo-server-express'
import jwt from 'jsonwebtoken'
import config from './config'
import { OAuth2Client } from 'google-auth-library'

const googleoauth = new OAuth2Client(config.auth.google.client_id)

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    hello: String
  }
`

// Provide resolver functions for your schema fields
const resolvers = {
  Query: {
    hello: (parent, args, thecontext) => {
      if (thecontext.user) console.log(`user ${JSON.stringify(thecontext.user)}`)
      else console.log('no user found')
      return 'Hello world!'
    },
  }
}

const context = async ({ req }) => {
  const strategy = req.headers.strategy
  if (strategy) {
    switch (strategy) {
      case 'google':
        console.log(`auth ${req.headers.authorization}`)
        return await googleoauth.getTokenInfo(req.headers.authorization.split(' ')[1])
          .then(user => {
            return { user }
          }).catch(err => {
            console.log(`got error ${err}`)
            throw new AuthenticationError('error getting token data from google')
            return
          })
      default:
        throw new AuthenticationError('invalid authentication strategy')
    }
  }
}

const server = new ApolloServer({
  typeDefs,
  resolvers,
  playground: true,
  tracing: true,
  context: context
})

const app = express()

app.use(cookieParser())

app.use(bodyParser.json())

app.use(bodyParser.urlencoded({
  extended: true
}))

app.use(cors({
  credentials: true,
  origin: true
}))

// Add headers
app.use(function (req, res, next) {

  // Website you wish to allow to connect
  res.setHeader('Access-Control-Allow-Origin', '*')

  // Request methods you wish to allow
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE')

  // Request headers you wish to allow
  res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type')

  // Pass to next layer of middleware
  next()
})

app.get('/user', (req, res) => {
  console.log('get user')
  if (req.headers.strategy === 'local') {
    jwt.verify(req.headers.authorization.split(' ')[1], config.auth.local.secret, (err, user) => {
      if (err) {
        res.status(401).json({
          message: err
        })
      } else {
        console.log(`got user ${JSON.stringify(user)}`)
        res.json(user)
      }
    })
  } else {
    res.status(401).json({
      message: 'no json found'
    })
  }
})

app.post('/logout', (req, res) => {
  console.log('log out')
  res.send('logged out I guess')
})
// put through graphql instead
app.post('/login', (req, res) => {
  console.log('login post request')
  console.log(req.body)
  // console.log(req)
  const username = req.body.username
  const password = req.body.password
  if (username && password) {
    jwt.sign({
      data: req.body
    }, config.auth.local.secret, {
        expiresIn: '1h'
      }, (err, token) => {
        if (err) {
          res.status(401).json({
            message: err
          })
        } else {
          console.log(`token ${token}`)
          res.json({
            token: token
          })
        }
      })
  } else {
    res.status(401).json({
      message: 'invalid username or password provided'
    })
  }
})

server.applyMiddleware({
  app,
  path: '/graphql',
  cors: false
})

const PORT = process.env.PORT || 8080

app.listen({ port: PORT }, () =>
  console.log(`Server listening on port ${PORT}, path ${server.graphqlPath} 🚀`)
)
