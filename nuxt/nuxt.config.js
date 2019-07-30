const pkg = require('./package')
require('dotenv').config()

const seodata = JSON.parse(process.env.SEOCONFIG)
const apiurl = process.env.APIURL
const ampurl = process.env.AMPURL
const recaptchasitekey = process.env.RECAPTCHASITEKEY

module.exports = {
  mode: 'spa',

  globalName: 'Joshua Schmidt',

  env: {
    seoconfig: process.env.SEOCONFIG,
    authconfig: process.env.AUTHCONFIG,
    apiurl: apiurl,
    ampurl: ampurl,
    shortlinkurl: process.env.SHORTLINKURL,
    recaptchasitekey: recaptchasitekey
  },

  /*
   ** Headers of the page
   */
  head: {
    titleTemplate: `%s - ${pkg.author}`,
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: pkg.description },
      // OpenGraph Data
      { property: 'og:title', content: 'Find Vericts' },
      { property: 'og:site_name', content: 'Find Vericts' },
      // The list of types is available here: http://ogp.me/#types
      { property: 'og:type', content: 'website' },
      {
        property: 'og:image',
        content: `${seodata.url}/opengraph.png`
      },
      { property: 'og:description', content: pkg.description },
      // Twitter card
      { name: 'twitter:card', content: 'summary' },
      {
        name: 'twitter:site',
        content: seodata.url
      },
      { name: 'twitter:title', content: 'Find Vericts' },
      {
        name: 'twitter:description',
        content: pkg.description
      },
      { name: 'twitter:creator', content: `@${seodata.twitterhandle}` },
      {
        name: 'twitter:image:src',
        content: `${seodata.url}/twitter.png`
      }
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }]
  },

  /*
   ** Router config
   */
  router: {},

  /*
   ** Customize the progress-bar color
   */
  loading: { color: '#fff' },

  /*
   ** Global CSS
   */
  css: [],

  /*
   ** Plugins to load before mounting the App
   */
  plugins: [
    { src: '~/plugins/font-awesome', ssr: false },
    { src: '~/plugins/vuelidate', ssr: false },
    { src: '~/plugins/vuex-persist', ssr: false },
    { src: '~/plugins/axios', ssr: false },
    { src: '~/plugins/toast', ssr: false },
    { src: '~/plugins/select', ssr: false },
    { src: '~/plugins/recaptcha', ssr: false },
    { src: '~/plugins/scroll-reveal', ssr: false }
  ],

  /*
   ** Nuxt.js modules
   */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
    // Doc: https://bootstrap-vue.js.org/docs/
    'bootstrap-vue/nuxt',
    '@nuxtjs/pwa',
    '@nuxtjs/style-resources',
    '@nuxtjs/dotenv',
    '@nuxtjs/sitemap',
    '@nuxtjs/google-analytics'
  ],

  /*
   ** google analytics config
   */
  googleAnalytics: {
    id: seodata.googleanalyticstrackingid
  },

  /*
   ** generate config
   */
  generate: {
    fallback: '404.html'
  },

  /*
   ** scss global config
   */
  styleResources: {
    scss: ['~assets/styles/global.scss']
  },

  /*
   ** Axios module configuration
   */
  axios: {
    // See https://github.com/nuxt-community/axios-module#options
    baseURL: apiurl
  },

  /*
   ** babel config
   */
  babel: {
    presets: ['es2015', 'stage-0'],
    plugins: [
      [
        'transform-runtime',
        {
          polyfill: true,
          regenerator: true
        }
      ]
    ]
  },

  /*
   ** Sitemap config
   */
  sitemap: {
    hostname: seodata.url,
    gzip: true,
    exclude: ['/admin/**']
  },

  extensions: ['js', 'ts'],

  /*
   ** Build configuration
   */
  build: {
    // put CSS in files instead of JS bundles
    extractCSS: true,
    /*
     ** You can extend webpack config here
     */
    extend(config, ctx) {
      // Run ESLint on save
      if (ctx.isDev && ctx.isClient) {
        config.module.rules.push({
          enforce: 'pre',
          test: /\.(js|vue)$/,
          loader: 'eslint-loader',
          exclude: /(node_modules)/
        })
      }
    }
  }
}
