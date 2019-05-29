import express = require('express')
import bodyParser = require('body-parser')
import { codes, adminconfig, mongoconfig } from './config'
import { initializeblogs } from './blogs'
import mongo from './mongo'

let db

mongo.then(client => {
  db = client.db(mongoconfig.dbname)
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
        _id: id
      }, {
          $set: {
            type: 'admin'
          }
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

adminApp.post('/initializeBlogs', (req, res) => {
  if (req.body.token === adminconfig.token) {
    initializeblogs().then(res1 => {
      res.json({
        message: res1
      }).status(codes.success)
    }).catch(err => {
      res.json({
        message: `blog init failed: ${err}`
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
