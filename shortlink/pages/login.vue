<template>
  <div></div>
</template>

<script lang="js">
import Vue from 'vue'
export default Vue.extend({
  name: 'Login',
  // @ts-ignore
  middleware: 'login',
  layout: 'redirect'
})
</script>

<style lang="scss"></style>
