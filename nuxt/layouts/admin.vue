<template>
  <div class="main-wrapper">
    <admin-navbar />
    <nuxt class="content" />
    <main-footer />
  </div>
</template>

<script lang="js">
import Vue from 'vue'
import AdminNavbar from '~/components/admin/Navbar.vue'
import MainFooter from '~/components/Footer.vue'
export default Vue.extend({
  name: 'Admin',
  // @ts-ignore
  middleware: 'admin',
  components: {
    AdminNavbar,
    MainFooter
  },
  // @ts-ignore
  head() {
    // @ts-ignore
    const seo = JSON.parse(process.env.seoconfig)
    const links = []
    const meta = []
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
      links,
      meta
    }
  }
})
</script>

<style lang="scss"></style>
