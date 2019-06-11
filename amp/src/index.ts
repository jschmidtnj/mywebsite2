import express = require('express')
import bodyParser = require('body-parser')
import { codes, config } from './config'
import blogtemplate from './blogtemplate'
import projecttemplate from './projecttemplate'
import axios from 'axios'
import * as cheerio from 'cheerio'
import { format } from 'date-fns'
import * as showdown from 'showdown'

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

ampApp.get('/projecttemplate', (req, res) => {
  res.send(projecttemplate)
    .status(codes.success)
})

ampApp.get('/blog/:id', (req, res) => {
  if (req.params.id) {
    const id = req.params.id
    let date
    try {
      const timestamp = id.toString().substring(0,8)
      date = format(parseInt(timestamp, 16) * 1000, 'M/D/YYYY')
    } catch(err) {
      res.json({
        code: codes.error,
        message: `error with timestamp parsing: ${err}`
      }).status(codes.error)
      return
    }
    axios.get(config.apiurl + '/graphql', {
      params: {
        query: `{post(type:"blog",id:"${id}"){title content author views}}`
      }
    }).then(res1 => {
      if (res1.status === 200) {
        if (res1.data) {
          if (res1.data.data && res1.data.data.post) {
            const blogdata = res1.data.data.post
            const $ = cheerio.load(blogtemplate)
            $('link[rel=canonical]').attr('href', config.websiteurl)
            $('#title').text(blogdata.title)
            $('#author').text(blogdata.author)
            const markdownconverter = new showdown.Converter()
            const blogcontenthtml = markdownconverter.makeHtml(blogdata.content)
            $('#content').text(blogcontenthtml)
            $('#views').text(blogdata.views)
            $('#date').text(date)
            $('#mainsite').attr('href', `${config.websiteurl}/blog/${id}`)
            $('title').text(blogdata.title)
            $('meta[name=description]').attr('content', `${blogdata.title} blog by ${blogdata.author}`)
            const html = $.html()
            res.send(html).status(codes.success)
          } else if (res1.data.errors) {
            res.json({
              code: codes.error,
              message: `found errors: ${JSON.stringify(res1.data.errors)}`
            }).status(codes.error)
          } else {
            res.json({
              code: codes.error,
              message: 'could not find data or errors'
            }).status(codes.error)
          }
        } else {
          res.json({
            code: codes.error,
            message: 'could not get data'
          }).status(codes.error)
        }
      } else {
        res.json({
          code: res1.status,
          message: `status code of ${res1.status}`
        }).status(res1.status)
      }
    }).catch(err => {
      res.json({
        code: codes.error,
        message: `got error ${err}`
      }).status(codes.error)
    })
  } else {
    res.json({
      code: codes.error,
      message: 'could not find id'
    }).status(codes.error)
  }
})

ampApp.get('/project/:id', (req, res) => {
  if (req.params.id) {
    const id = req.params.id
    let date
    try {
      const timestamp = id.toString().substring(0,8)
      date = format(parseInt(timestamp, 16) * 1000, 'M/D/YYYY')
    } catch(err) {
      res.json({
        code: codes.error,
        message: `error with timestamp parsing: ${err}`
      }).status(codes.error)
      return
    }
    axios.get(config.apiurl + '/graphql', {
      params: {
        query: `{post(type:"project",id:"${id}"){title content author views}}`
      }
    }).then(res1 => {
      if (res1.status === 200) {
        if (res1.data) {
          if (res1.data.data && res1.data.data.post) {
            const projectdata = res1.data.data.post
            const $ = cheerio.load(blogtemplate)
            $('link[rel=canonical]').attr('href', config.websiteurl)
            $('#title').text(projectdata.title)
            $('#author').text(projectdata.author)
            const markdownconverter = new showdown.Converter()
            const projectcontenthtml = markdownconverter.makeHtml(projectdata.content)
            $('#content').text(projectcontenthtml)
            $('#views').text(projectdata.views)
            $('#date').text(date)
            $('#mainsite').attr('href', `${config.websiteurl}/project/${id}`)
            $('title').text(projectdata.title)
            $('meta[name=description]').attr('content', `${projectdata.title} project by ${projectdata.author}`)
            const html = $.html()
            res.send(html).status(codes.success)
          } else if (res1.data.errors) {
            res.json({
              code: codes.error,
              message: `found errors: ${JSON.stringify(res1.data.errors)}`
            }).status(codes.error)
          } else {
            res.json({
              code: codes.error,
              message: 'could not find data or errors'
            }).status(codes.error)
          }
        } else {
          res.json({
            code: codes.error,
            message: 'could not get data'
          }).status(codes.error)
        }
      } else {
        res.json({
          code: res1.status,
          message: `status code of ${res1.status}`
        }).status(res1.status)
      }
    }).catch(err => {
      res.json({
        code: codes.error,
        message: `got error ${err}`
      }).status(codes.error)
    })
  } else {
    res.json({
      code: codes.error,
      message: 'could not find id'
    }).status(codes.error)
  }
})

const PORT = process.env.PORT || config.port

ampApp.listen(PORT, () => console.log(`amp app is listening on port ${PORT} 🚀`))
