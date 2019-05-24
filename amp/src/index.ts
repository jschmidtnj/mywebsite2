import express = require('express')
import bodyParser = require('body-parser')
import { codes, config } from './config'
import axios from 'axios'

const cronApp = express()

cronApp.use(
  bodyParser.urlencoded({
    extended: false
  })
)

cronApp.use(bodyParser.json())

cronApp.get('/hello', (req, res) => {
  res
    .json({
      message: `Hello!`,
      code: codes.success
    })
    .status(codes.success)
})

const PORT = process.env.PORT || config.port

cronApp.listen(PORT, () => console.log(`cron app is listening on port ${PORT} 🚀`))
