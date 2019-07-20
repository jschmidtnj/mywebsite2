<template>
  <div>
    <admin-navbar v-if="$auth.user && $auth.user.type === 'admin'" />
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
  }
})
</script>

<style lang="scss"></style>
