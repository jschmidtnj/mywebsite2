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
  }
})
</script>

<style lang="scss"></style>
