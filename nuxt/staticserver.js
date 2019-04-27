const connect = require('connect')
const serveStatic = require('serve-static')

const PORT = process.env.PORT || 8080

/* eslint-disable */
connect()
  .use(serveStatic('dist'))
  .listen(PORT, () => {
    console.log(`cron app is listening on port ${PORT} 🚀`)
  })
