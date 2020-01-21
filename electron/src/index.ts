import nativefier from 'nativefier'
import { Storage } from '@google-cloud/storage'
import * as archiver from 'archiver'
import config from './config'

const options = {
  name: config.name,
  targetUrl: config.url,
  version: '1.0.0',
  out: './dist/',
  overwrite: true,
  icon: './assets/icon.ico',
  counter: false,
  bounce: false,
  width: 1280,
  height: 800,
  showMenuBar: false,
  fastQuit: false,
  ignoreCertificate: false,
  ignoreGpuBlacklist: false,
  enableEs3Apis: false,
  insecure: false,
  honest: false,
  zoom: 1.0,
  singleInstance: false,
  clearCache: false,
  fileDownloadOptions: {
    saveAs: true // always show "Save As" dialog
  },
  processEnvs: {}
}

const storage = new Storage({
  projectId: config.storageConfig.project_id,
  credentials: config.storageConfig
})
const filepath = `releases/${process.platform}`
const writestream = storage.bucket(config.storageBucketName).file(filepath).createWriteStream({
  metadata: {
    resumable: false,
    contentType: 'application/zip',
    predefinedAcl: 'publicRead',
    metadata: {} // any extra data can be added here
  }
})
writestream.on('finish', () => {
  console.log(`uploaded ${filepath}`)
})
writestream.on('error', err => {
  console.error(err)
})

const archive = archiver('zip', {
  zlib: {
    level: 9 // compression level
  }
})

archive.on('close', () => {
  console.log(`${archive.pointer()} total bytes`)
  writestream.end()
})

archive.on('warning', function(err) {
  if (err.code === 'ENOENT') {
    console.log(err)
  } else {
    console.error(err)
    throw err
  }
})

archive.on('error', (err) => {
  console.error(err)
  writestream.end()
  throw err
})

archive.pipe(writestream)

nativefier(options, (err, appPath) => {
  if (err) {
    console.error(err)
    return
  }
  console.log('App has been nativefied to', appPath)
  archive.directory(appPath, false)
  archive.finalize().then(() => {
    console.log('all finished')
  }).catch(err1 => {
    console.error(err1)
  })
})
