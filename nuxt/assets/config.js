export const androidAppURL = 'https://play.google.com/apps/testing/com.joshuaschmidt.joshuaschmidt'

export const cloudStorageURLs = {
  posts: 'https://cdn.joshuaschmidt.tech',
  desktop: 'https://storage.googleapis.com/mywebsite2desktop'
}

export const codes = {
  error: 400,
  success: 200,
  warning: 300,
  unauthorized: 403
}

export const toasts = {
  position: 'top-right',
  duration: 2000,
  theme: 'bubble'
}

export const validTypes = ['blog', 'project']

export const regex = {
  password: /^$|^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,}$/,
  hexcode: /(^#[0-9A-F]{6}$)|(^#[0-9A-F]{3}$)/i
}

export const options = {
  categoryOptions: ['technology', 'webdesign'],
  tagOptions: ['vue', 'nuxt']
}

export const defaultColor = '#194d332B'

export const staticstorageindexes = {
  blogfiles: 'blogfiles',
  projectfiles: 'projectfiles',
  placeholder: 'placeholder'
}

export const paths = {
  placeholder: '/placeholder',
  original: '/original',
  blur: '/blur'
}

export const validimages = [
  'image/jpeg',
  'image/png',
  'image/svg+xml'
]

export const validfiles = [
  'image/jpeg',
	'image/png',
	'image/svg+xml',
	'image/gif',
  'video/mpeg',
  'video/mp4',
	'video/webm',
	'video/x-msvideo',
	'application/pdf',
	'text/plain',
	'application/zip',
	'text/csv',
	'application/json',
	'application/ld+json',
	'application/vnd.ms-powerpoint',
	'application/vnd.openxmlformats-officedocument.presentationml.presentation',
	'application/msword',
	'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
]
