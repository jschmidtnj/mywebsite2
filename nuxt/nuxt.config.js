const pkg = require('./package')
require('dotenv').config()

const seodata = JSON.parse(process.env.SEOCONFIG)
const authdata = JSON.parse(process.env.AUTHCONFIG)

module.exports = {
  mode: 'universal',

  env: {
    seoconfig: process.env.SEOCONFIG,
    authconfig: process.env.AUTHCONFIG
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
    { src: '~/plugins/axios', ssr: false }
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
    '@nuxtjs/google-analytics',
    '@nuxtjs/auth'
  ],

  /*
   ** auth config
   */
  auth: {
    strategies: {
      local: {
        endpoints: {
          login: {
            url: '/login',
            method: 'post',
            propertyName: 'token'
          },
          logout: {
            url: '/logout',
            method: 'post'
          },
          user: {
            url: '/user',
            method: 'get',
            propertyName: 'user'
          }
        }
      },
      google: {
        client_id: authdata.google.client_id
      }
    },
    redirect: {
      callback: '/callback',
      logout: '/login',
      login: '/login',
      home: '/account'
    },
    resetOnError: true
  },

  /*
   ** google analytics config
   */
  googleAnalytics: {
    id: seodata.googleanalyticstrackingid
  },

  /*
   ** generate 404 page
   */
  generate: {
    fallback: '404.html'
  },

  /*
   ** scss global config
   */
  styleResources: {
    scss: ['~assets/styles/bootstrap.scss']
  },

  /*
   ** Axios module configuration
   */
  axios: {
    // See https://github.com/nuxt-community/axios-module#options
    baseURL: 'http://localhost:8080'
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
