/**
 * sample config for init scripts - change values and rename config.ts, place in this directory
 */

export const elasticuri = 'elasticuri'

export const adminconfig = {
  port: 8080,
  token: '123',
  saltrounds: 15
}

export const codes = {
  error: 400,
  success: 200,
  warning: 300,
  unauthorized: 403
}

export const mongoconfig = {
  uri: 'mongouri',
  dbname: 'testing'
}
