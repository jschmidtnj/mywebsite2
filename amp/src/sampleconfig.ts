/**
 * sample config for amp - change values and rename config.ts, place in this directory
 */

export const config = {
  port: 8080
}

export const codes = {
  error: 400,
  success: 200,
  warning: 300,
  unauthorized: 403
}

export const firebase = {
  databaseURL: 'https://project.firebaseio.com',
  storageBucket: 'project.appspot.com',
  projectId: 'project',
  credentials: {
    type: 'service_account',
    project_id: 'project',
    private_key_id: 'id',
    private_key: 'key',
    client_email: 'email@project.iam.gserviceaccount.com',
    client_id: 'id',
    auth_uri: 'https://accounts.google.com/o/oauth2/auth',
    token_uri: 'https://oauth2.googleapis.com/token',
    auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
    client_x509_cert_url: 'url'
  }
}
