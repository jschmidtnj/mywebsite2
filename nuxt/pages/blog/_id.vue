<template>
  <div>
    <b-container v-if="blog">
      <h1>{{ blog.title }}</h1>
      <p>{{ blog.author }}</p>
      <p v-if="blog.date">{{ formatDate(blog.date, 'M/D/YYYY') }}</p>
      <p>{{ blog.views }}</p>
      <p>{{ blog.content }}</p>
    </b-container>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { format } from 'date-fns'
export default Vue.extend({
  name: 'Blog',
  /* eslint-disable */
  // @ts-ignore
  asyncData(context) {
    // @ts-ignore
    const handleErrors = () => {
      return {
        id: null,
        blog: {}
      }
    }
    if (context.params && context.params.id) {
      const id = context.params.id
      return context.app.$axios
        .get('/graphql', {
          params: {
            query: `{blog(id:"${id}"){title content id author views}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.blog) {
                const blog = res.data.data.blog
                console.log(res.data.data.blog)
                return {
                  id: id,
                  blog: blog
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
    const title = this.blog ? this.blog.title : 'Blog'
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
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    }
  }
})
</script>

<style lang="scss"></style>
