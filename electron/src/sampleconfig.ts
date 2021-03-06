export default {
  storageConfig: {
    type: 'service_account',
    project_id: 'id',
    private_key_id: 'id',
    private_key: 'key',
    client_email: 'account@name.iam.gserviceaccount.com',
    client_id: 'id',
    auth_uri: 'https://accounts.google.com/o/oauth2/auth',
    token_uri: 'https://oauth2.googleapis.com/token',
    auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
    client_x509_cert_url: 'https://www.googleapis.com/robot/v1/metadata/x509/projectid.iam.gserviceaccount.com'
  },
  storageBucketName: 'storagename',
  name: 'mywebsite',
  url: 'https://example.com'
}
