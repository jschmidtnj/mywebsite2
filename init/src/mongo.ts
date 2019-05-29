import * as mongodb from 'mongodb'
import { mongoconfig } from './config'

const MongoClient = mongodb.MongoClient

const client = new MongoClient(mongoconfig.uri, {
  useNewUrlParser: true
})

export default client.connect()
