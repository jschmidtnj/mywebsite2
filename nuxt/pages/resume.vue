<template>
  <b-container>
    <v-pdf :src="pdf" :resize="true" :annotation="true"></v-pdf>
    <a :href="pdf" download="joshua_schmidt_resume" class="mt-2">Download</a>
  </b-container>
</template>

<script lang="js">
import Vue from 'vue'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
const pdf = process.env.pdf
export default Vue.extend({
  name: 'Resume',
  data() {
    return {
      pdf
    }
  },
  // @ts-ignore
  head() {
    const title = 'Resume'
    const description = 'view my resume'
    const image = `${seo.url}/icon.png`
    return {
      title,
      meta: [
        { property: 'og:title', content: title },
        { property: 'og:description', content: description },
        {
          property: 'og:image',
          content: image
        },
        { name: 'twitter:title', content: title },
        {
          name: 'twitter:description',
          content: description
        },
        {
          name: 'twitter:image',
          content: image
        },
        { hid: 'description', name: 'description', content: description }
      ]
    }
  }
})
</script>

<style lang="scss"></style>
