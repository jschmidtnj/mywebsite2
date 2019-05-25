import express = require('express')
import bodyParser = require('body-parser')
import { codes, config } from './config'
import blogtemplate from './blogtemplate'
import axios from 'axios'

const ampApp = express()

ampApp.use(
  bodyParser.urlencoded({
    extended: false
  })
)

ampApp.use(bodyParser.json())

ampApp.get('/hello', (req, res) => {
  res
    .json({
      message: `Hello!`,
      code: codes.success
    })
    .status(codes.success)
})

ampApp.get('/blogtemplate', (req, res) => {
  res.send(blogtemplate)
    .status(codes.success)
})

const PORT = process.env.PORT || config.port

ampApp.listen(PORT, () => console.log(`amp app is listening on port ${PORT} 🚀`))
