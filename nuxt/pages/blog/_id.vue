<template>
  <div>
    <b-container v-if="post">
      <h1>{{ post.title }}</h1>
      <p>{{ post.author }}</p>
      <p v-if="post.id">{{ formatDate(mongoidToDate(post.id), 'M/D/YYYY') }}</p>
      <p>{{ post.views }}</p>
      <vue-markdown
        :source="post.content"
        class="mb-4 markdown"
        @rendered="updateCodeHighlighting"
      />
    </b-container>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { format } from 'date-fns'
import VueMarkdown from 'vue-markdown'
import Prism from 'prismjs'
export default Vue.extend({
  name: 'Post',
  components: {
    VueMarkdown
  },
  /* eslint-disable */
  // @ts-ignore
  asyncData(context) {
    // @ts-ignore
    const handleErrors = () => {
      return {
        id: null,
        post: {}
      }
    }
    if (context.params && context.params.id) {
      const id = context.params.id
      return context.app.$axios
        .get('/graphql', {
          params: {
            query: `{post(type:"blog",id:"${id}"){title content id author views}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const post = res.data.data.post
                console.log(res.data.data.post)
                return {
                  id: id,
                  post: post
                }
              } else if (res.data.errors) {
                console.log(`found errors: ${JSON.stringify(res.data.errors)}`)
                return handleErrors()
              } else {
                console.log('could not find data or errors')
                return handleErrors()
              }
            } else {
              console.log('could not get data')
              return handleErrors()
            }
          } else {
            console.log(`status code of ${res.status}`)
            return handleErrors()
          }
        })
        .catch(err => {
          console.error(`got error: ${err}`)
          return handleErrors()
        })
    } else {
      console.log('could not find id in params')
      return handleErrors()
    }
  },
  // @ts-ignore
  head() {
    const title = this.post ? this.post.title : 'Post'
    // @ts-ignore
    if (!(process.env.seoconfig && process.env.ampurl)) {
      return {
        title: title,
        link: []
      }
    }
    // @ts-ignore
    const seo = JSON.parse(process.env.seoconfig)
    // @ts-ignore
    const ampurl = process.env.ampurl
    return {
      title: title,
      link: [
        {
          rel: 'canonical',
          href: `${seo.url}/blog`
        },
        {
          rel: 'amphtml',
          href: `${ampurl}/blog/${this.$route.params.id}`
        }
      ]
    }
  },
  methods: {
    updateCodeHighlighting() {
      this.$nextTick(() => {
        Prism.highlightAll()
      })
    },
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    },
    mongoidToDate(id) {
      return parseInt(id.substring(0,8), 16) * 1000
    }
  }
})
</script>

<style lang="scss">
@import '~/node_modules/prismjs/themes/prism.css';
</style>
