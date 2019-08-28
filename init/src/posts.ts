import elasticsearch = require('elasticsearch')
import { elasticuri } from './config'

/**
 * post functions - initialize posts
 */

const postmappings = {
  properties: {
    title: {
      type: 'keyword'
    },
    caption: {
      type: 'keyword'
    },
    author: {
      type: 'keyword'
    },
    color: {
      type: 'text'
    },
    tags: {
      type: 'keyword'
    },
    categories: {
      type: 'keyword'
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
    },
    heroimage: {
      type: 'nested'
    },
    tileimage: {
      type: 'nested'
    },
    files: {
      type: 'nested'
    },
    comments: {
      type: 'nested'
    }
  }
}

const postindexsettings = {
  number_of_shards: 2,
  number_of_replicas: 1
}

const writeclient = new elasticsearch.Client({
  host: elasticuri
})

export const initializeposts = (postindexname, postdoctype) => {
  return new Promise((resolve, reject) => {
    writeclient
      .ping({
        requestTimeout: 1000
      })
      .then(() => {
        console.log(`able to ping writeclient`)
        writeclient.indices
          .putSettings({
            index: postindexname,
            body: {
              index: {
                number_of_replicas: 0,
                number_of_shards: 1
              }
            }
          })
          .then(res0 => {
            console.log(
              `deleted all shards in ${postindexname}: ${JSON.stringify(
                res0
              )}`
            )
          })
          .catch(err => {
            const errmessage = `error deleting shards in index ${postindexname}: ${err}`
            console.log(errmessage)
          })
          .then(() => {
            writeclient.indices
              .delete({
                index: postindexname
              })
              .then(res1 => {
                console.log(
                  `deleted index ${postindexname}: ${JSON.stringify(res1)}`
                )
              })
              .catch(err => {
                const errmessage = `error deleting index ${postindexname}: ${err}`
                console.log(errmessage)
              })
              .then(() => {
                return writeclient.indices
                  .exists({
                    index: postindexname
                  })
                  .then(res2 => {
                    if (res2) {
                      console.log(`index ${postindexname} exists still`)
                    } else {
                      return writeclient.indices
                        .create({
                          index: postindexname,
                          body: {
                            settings: postindexsettings
                          }
                        })
                        .then(res3 => {
                          console.log(`added index ${postindexname}: ${res3}`)
                          return writeclient.indices
                            .getMapping()
                            .then(res4 => {
                              if (
                                Object.keys(res4[postindexname].mappings)
                                  .length === 0
                              ) {
                                console.log(
                                  `${postindexname}: no mappings :)`
                                )
                                return writeclient.indices
                                  .putMapping({
                                    index: postindexname,
                                    type: postdoctype,
                                    body: postmappings,
                                    include_type_name: true
                                  })
                                  .then(res5 => {
                                    console.log(
                                      `initialized ${postindexname}: ${JSON.stringify(
                                        res5
                                      )}`
                                    )
                                    resolve(`finished initializing elasticsearch`)
                                  })
                                  .catch(err => {
                                    const errmessage = `could not create mapping for ${postindexname}: ${err}`
                                    console.log(errmessage)
                                    reject(errmessage)
                                  })
                              } else {
                                const errmessage = `${postindexname} already has mappings :(`
                                console.log(errmessage)
                                reject(errmessage)
                              }
                            })
                            .catch(err => {
                              const errmessage = `could not get mappings for ${postindexname}: ${err}`
                              console.log(errmessage)
                              reject(errmessage)
                            })
                        })

                        .catch(err => {
                          const errmessage = `error adding index ${postindexname}: ${err}`
                          console.log(errmessage)
                          reject(errmessage)
                        })
                    }
                  })
                  .catch(err => {
                    const errmessage = `error checking if index ${postindexname} exists: ${err}`
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
