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
  // @ts-ignore
  middleware: 'checkloggedin',
  components: {
    Navbar
  }
})
</script>

<style lang="scss">
@import '~/node_modules/progressive-image.js/dist/progressive-image.css';
</style>
