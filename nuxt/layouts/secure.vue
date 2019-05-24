<template>
  <div>
    <nuxt />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'Secure',
  // @ts-ignore
  middleware: 'auth'
})
</script>

<style lang="scss"></style>
