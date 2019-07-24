<template>
  <div>
    <admin-navbar v-if="admin" />
    <nuxt />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import AdminNavbar from '~/components/admin/Navbar.vue'
export default Vue.extend({
  name: 'Secure',
  // @ts-ignore
  middleware: 'auth',
  components: {
    AdminNavbar
  },
  computed: {
    admin() {
      return (
        this.$store.state.auth &&
        this.$store.state.auth.user &&
        this.$store.state.auth.user.type === 'admin'
      )
    }
  }
})
</script>

<style lang="scss"></style>
