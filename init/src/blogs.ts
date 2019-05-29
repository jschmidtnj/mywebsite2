import elasticsearch = require('elasticsearch')
import { elasticuri, mongoconfig } from './config'
import mongo from './mongo'

let db

mongo.then(client => {
  db = client.db(mongoconfig.dbname)
}).catch(err => {
  console.error(`got error connecting to mongo: ${err}`)
})

/**
 * blog functions - editing / searching blogs etc.
 */

const blogindexname = 'blogs'
const blogdoctype = 'blog'
const blogmappings = {
  properties: {
    title: {
      type: 'text'
    },
    author: {
      type: 'text'
    },
    content: {
      type: 'text'
    },
    views: {
      type: 'integer'
    },
    date: {
      type: 'date',
      format: 'epoch_millis'
    }
  }
}

const blogindexsettings = {
  number_of_shards: 2,
  number_of_replicas: 1
}

const writeclient = new elasticsearch.Client({
  host: elasticuri
})

export const initializeblogs = () => {
  return new Promise((resolve, reject) => {
    writeclient
      .ping({
        requestTimeout: 1000
      })
      .then(() => {
        console.log(`able to ping writeclient`)
        writeclient.indices
          .putSettings({
            index: blogindexname,
            body: {
              index: {
                number_of_replicas: 0,
                number_of_shards: 0
              }
            }
          })
          .then(res0 => {
            console.log(
              `deleted all shards in ${blogindexname}: ${JSON.stringify(
                res0
              )}`
            )
          })
          .catch(err => {
            const errmessage = `error deleting shards in index ${blogindexname}: ${err}`
            console.log(errmessage)
          })
          .then(() => {
            writeclient.indices
              .delete({
                index: blogindexname
              })
              .then(res1 => {
                console.log(
                  `deleted index ${blogindexname}: ${JSON.stringify(res1)}`
                )
              })
              .catch(err => {
                const errmessage = `error deleting index ${blogindexname}: ${err}`
                console.log(errmessage)
              })
              .then(() => {
                return writeclient.indices
                  .exists({
                    index: blogindexname
                  })
                  .then(res2 => {
                    if (res2) {
                      console.log(`index ${blogindexname} exists still`)
                    } else {
                      return writeclient.indices
                        .create({
                          index: blogindexname,
                          body: {
                            settings: blogindexsettings
                          }
                        })
                        .then(res3 => {
                          console.log(`added index ${blogindexname}: ${res3}`)
                          return writeclient.indices
                            .getMapping()
                            .then(res4 => {
                              if (
                                Object.keys(res4[blogindexname].mappings)
                                  .length === 0
                              ) {
                                console.log(
                                  `${blogindexname}: no mappings :)`
                                )
                                return writeclient.indices
                                  .putMapping({
                                    index: blogindexname,
                                    type: blogdoctype,
                                    body: blogmappings
                                  })
                                  .then(res5 => {
                                    console.log(
                                      `initialized ${blogindexname}: ${JSON.stringify(
                                        res5
                                      )}`
                                    )
                                    resolve(`finished initializing elasticsearch`)
                                    // get from mongodb
                                    db
                                      .collection('blogs')
                                      .find({})
                                      .toArray()
                                      .then(docs => {
                                        if (docs.length === 0) {
                                          resolve(`finished initializing elasticsearch: ${docs.length} blogs`)
                                        } else {
                                          let responsecount = 0
                                          const delay = 100
                                          let iterationcount = -1
                                          docs.forEach(doc => {
                                            const id = doc._id
                                            delete doc._id
                                            doc.date = new Date(doc._id.getTimestamp()).getTime()
                                            iterationcount++
                                            setTimeout(() => {
                                              console.log(`updating blog ${id}`)
                                              writeclient
                                                .index({
                                                  index: blogindexname,
                                                  type: blogdoctype,
                                                  id: id,
                                                  body: doc
                                                })
                                                .then(res6 => {
                                                  responsecount++
                                                  console.log(responsecount, docs.length)
                                                  console.log(
                                                    `created / updated ${
                                                    doc.title
                                                    }: ${JSON.stringify(res6)}`
                                                  )
                                                  if (responsecount === docs.length) {
                                                    resolve(
                                                      `finished updating everything`
                                                    )
                                                  }
                                                })
                                                .catch(err => {
                                                  const errmessage = `could not create / update ${
                                                    doc.title
                                                    }: ${err}`
                                                  console.log(errmessage)
                                                  reject(errmessage)
                                                })
                                            }, delay * iterationcount)
                                          })
                                        }
                                      }).catch(err => reject(err))
                                  })
                                  .catch(err => {
                                    const errmessage = `could not create mapping for ${blogindexname}: ${err}`
                                    console.log(errmessage)
                                    reject(errmessage)
                                  })
                              } else {
                                const errmessage = `${blogindexname} already has mappings :(`
                                console.log(errmessage)
                                reject(errmessage)
                              }
                            })
                            .catch(err => {
                              const errmessage = `could not get mappings for ${blogindexname}: ${err}`
                              console.log(errmessage)
                              reject(errmessage)
                            })
                        })

                        .catch(err => {
                          const errmessage = `error adding index ${blogindexname}: ${err}`
                          console.log(errmessage)
                          reject(errmessage)
                        })
                    }
                  })
                  .catch(err => {
                    const errmessage = `error checking if index ${blogindexname} exists: ${err}`
                    console.log(errmessage)
                    reject(errmessage)
                  })
              })
          })
      })
      .catch(err => {
        const errmessage = `unable to ping writeclient`
        console.log(errmessage)
        reject(err)
      })
  })
}
