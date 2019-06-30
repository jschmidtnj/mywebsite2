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
    <loading v-else />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { format } from 'date-fns'
import VueMarkdown from 'vue-markdown'
import Loading from '~/components/Loading.vue'
import Prism from 'prismjs'
import { validTypes } from '~/assets/config'
// @ts-ignore
const ampurl = process.env.ampurl
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  name: 'Post',
  components: {
    VueMarkdown,
    Loading
  },
  props: {
    type: {
      type: String,
      default: null,
      required: true,
      validator: val => validTypes.includes(String(val))
    }
  },
  data() {
    return {
      id: null,
      post: null
    }
  },
  /* eslint-disable */
  mounted() {
    if (this.$route.query && this.$route.query.id) {
      this.id = this.$route.query.id
      // update document amphtml for spa
      // @ts-ignore
      document.head.querySelector("link[rel='amphtml']").href = `${ampurl}/blog/${this.id}`
      // @ts-ignore
      document.head.querySelector("link[rel='canonical']").href = `${seo.url}/blog?id=${this.$route.query.id}`
      this.$axios
        .get('/graphql', {
          params: {
            query: `{post(type:"${this.type}",id:"${this.id}"){title content id author views}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const post = res.data.data.post
                post.content = decodeURIComponent(post.content)
                this.post = post
                // update title for spa
                document.title = this.post.title
              } else if (res.data.errors) {
                this.$toasted.global.error({
                  message: `found errors: ${JSON.stringify(res.data.errors)}`
                })
              } else {
                this.$toasted.global.error({
                  message: 'could not find data or errors'
                })
              }
            } else {
              this.$toasted.global.error({
                message: 'could not get data'
              })
            }
          } else {
            this.$toasted.global.error({
              message: `status code of ${res.status}`
            })
          }
        })
        .catch(err => {
          this.$toasted.global.error({
            message: `got error: ${err}`
          })
        })
    } else {
      this.$toasted.global.error({
        message: 'could not find id in params'
      })
    }
  },
  // @ts-ignore
  head() {
    const title = this.post ? this.post.title : validTypes.includes(this.type) ? this.type : 'Post'
    return {
      title: title,
      link: [
        {
          rel: 'canonical',
          href: `${seo.url}/blog?id=${this.$route.query.id}`
        },
        {
          rel: 'amphtml',
          href: `${ampurl}/blog/${this.$route.query.id}`
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
