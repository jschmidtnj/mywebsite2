import express = require('express')
import bodyParser = require('body-parser')
import { codes, adminconfig, mongoconfig } from './config'
import { initializeposts } from './posts'
import * as mongodb from 'mongodb'

// indexname must match mongodb name
const blogIndexName = 'blogs'
const blogDocType = 'blog'
const projectIndexName = 'projects'
const projectDocType = 'project'

const MongoClient = mongodb.MongoClient
const ObjectID = mongodb.ObjectID

const client = new MongoClient(mongoconfig.uri, {
  useNewUrlParser: true,
  useUnifiedTopology: true
})

let db

client.connect().then(theclient => {
  db = theclient.db(mongoconfig.dbname)
}).catch(err => {
  console.error(`got error connecting to mongo: ${err}`)
})

const adminApp = express()

adminApp.use(
  bodyParser.urlencoded({
    extended: false
  })
)

adminApp.use(bodyParser.json())

adminApp.get('/hello', (req, res) => {
  res
    .json({
      message: `Hello!`,
      code: codes.success
    })
    .status(codes.success)
})

adminApp.post('/addAdmin', (req, res) => {
  if (req.body.token === adminconfig.token) {
    const id = req.body.id
    if (id) {
      db.collection('users').updateOne({
        _id: new ObjectID(id)
      },
        {
          $set: {
            type: 'admin'
          }
        },
        {
          upsert: true
        }).then(res1 => {
          res.json({
            message: `updated user ${id} to admin`
          })
            .status(codes.success)
        }).catch(err => {
          res.json({
            message: `error updating to admin: ${err}`
          })
            .status(codes.error)
        })
    } else {
      res.json({
        message: `no id provided`
      })
        .status(codes.error)
    }
  } else {
    res.json({
      message: `Invalid admin token`
    })
      .status(codes.unauthorized)
  }
})

adminApp.post('/initializePosts', (req, res) => {
  if (req.body.token === adminconfig.token) {
    initializeposts(blogIndexName, blogDocType).then(res1 => {
      initializeposts(projectIndexName, projectDocType).then(res2 => {
        res.json({
          message: `res1: ${res1}, res2: ${res2}`
        }).status(codes.success)
      }).catch(err => {
        res.json({
          message: `post init failed: ${err}`
        })
          .status(codes.error)
      })
    }).catch(err => {
      res.json({
        message: `post init failed: ${err}`
      })
        .status(codes.error)
    })
  } else {
    res.json({
      message: `Invalid admin token`
    })
      .status(codes.unauthorized)
  }
})

const PORT = process.env.PORT || adminconfig.port

adminApp.listen(PORT, () => console.log(`admin app is listening on port ${PORT} 🚀`))
