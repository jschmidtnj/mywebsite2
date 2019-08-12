<template>
  <div class="main-wrapper">
    <admin-navbar v-if="admin" />
    <nuxt class="content" />
    <main-footer />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import AdminNavbar from '~/components/admin/Navbar.vue'
import MainFooter from '~/components/Footer.vue'
export default Vue.extend({
  name: 'Secure',
  // @ts-ignore
  middleware: 'auth',
  components: {
    AdminNavbar,
    MainFooter
  },
  computed: {
    admin() {
      return (
        this.$store.state.auth &&
        this.$store.state.auth.user &&
        this.$store.state.auth.user.type === 'admin'
      )
    }
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
