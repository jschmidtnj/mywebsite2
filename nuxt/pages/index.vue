<template>
  <b-container class="my-4">
    <div>
      <!--b-jumbotron
        header="Joshua Schmidt"
        lead="Programmer, Computer Engineer, Developer"
        class="mt-4"
      >
        <b-button variant="primary" href="/about">About</b-button>
      </b-jumbotron-->
      <h3>Blogs</h3>
      <tile-rows type="blog" />
      <h3>Projects</h3>
      <tile-rows type="project" />
    </div>
  </b-container>
</template>

<script lang="js">
import Vue from 'vue'
import TileRows from '~/components/TileRows.vue'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  components: {
    TileRows
  },
  data() {
    return {}
  },
  mounted() {},
  // @ts-ignore
  head() {
    const title = 'Home'
    const description =
      'custom blog and project descriptions for anything in tech and beyond'
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
      ],
      __dangerouslyDisableSanitizers: ['script'],
      script: [
        {
          innerHTML: JSON.stringify({
            '@context': 'http://schema.org',
            '@type': 'WebPage',
            name: title,
            description
          }),
          type: 'application/ld+json'
        }
      ]
    }
  }
})
</script>

<style lang="scss"></style>
