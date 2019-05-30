<template>
  <b-container>
    <h1>{{ blog.title }}</h1>
    <p>{{ blog.author }}</p>
    <p v-if="blog.date">{{ formatDate(blog.date, 'M/D/YYYY') }}</p>
    <p>{{ blog.views }}</p>
    <p>{{ blog.content }}</p>
  </b-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { format } from 'date-fns'
export default Vue.extend({
  name: 'Blog',
  data() {
    return {
      id: null,
      blog: {}
    }
  },
  mounted() {
    /* eslint-disable */
    if (this.$route.query && this.$route.query.id) {
      this.id = this.$route.query.id
      this.$axios.get('/graphql', {
        params: {
          query: `{blog(id:"${this.id}"){title content id author views}}`
        }
      }).then(res => {
        if (res.status === 200) {
          if (res.data) {
            if (res.data.data && res.data.data.blog) {
              this.blog = res.data.data.blog
              console.log(res.data.data.blog)
            } else if (res.data.errors) {
              console.log(`found errors: ${res.data.errors}`)
            } else {
              console.log('could not find data or errors')
            }
          } else {
            console.log('could not get data')
          }
        } else {
          console.log(`status code of ${res.status}`)
        }
      }).catch(err => {
        console.error(`got error: ${err}`)
      })
    } else {
      console.log('could not find id in query')
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
