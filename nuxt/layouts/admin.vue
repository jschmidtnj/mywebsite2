<template>
  <div>
    <admin-navbar />
    <nuxt />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import AdminNavbar from '~/components/admin/Navbar.vue'
export default Vue.extend({
  name: 'Admin',
  // @ts-ignore
  middleware: 'admin',
  components: {
    AdminNavbar
  }
})
</script>

<style lang="scss"></style>
