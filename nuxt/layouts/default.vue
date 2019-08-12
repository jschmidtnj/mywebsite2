<template>
  <div>
    <navbar />
    <nuxt />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Navbar from '~/components/Navbar.vue'

export default Vue.extend({
  name: 'Default',
  components: {
    Navbar
  },
  // @ts-ignore
  head() {
    // @ts-ignore
    const seo = JSON.parse(process.env.seoconfig)
    const links: any = []
    const meta: any = []
    if (seo) {
      const canonical = `${seo.url}/${this.$route.path}`
      links.push({
        rel: 'canonical',
        href: canonical
      })
      meta.push({
        property: 'og:url',
        content: canonical
      })
    }
    return {
      links: links,
      meta: meta
    }
  }
})
</script>

<style lang="scss"></style>
