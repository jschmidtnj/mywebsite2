<template>
  <b-container>
    <v-pdf :src="pdf" class="pdf"></v-pdf>
  </b-container>
</template>

<script lang="js">
import Vue from 'vue'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
const pdf = process.env.pdf
export default Vue.extend({
  name: 'Resume',
  // @ts-ignore
  head() {
    const title = 'Resume'
    const description = 'view my resume'
    const image = `${seo.url}/icon.png`
    return {
      title: title,
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
  },
  data() {
    return {
      pdf: pdf
    }
  }
})
</script>

<style lang="scss">
.pdf {
  display: inline-block;
  width: 100%;
}
</style>
