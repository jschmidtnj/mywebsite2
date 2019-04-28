<template>
  <div>
    <p>Secure Page</p>
    <b-btn block @click="$auth.logout()">Logout</b-btn>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  // @ts-ignore
  layout: 'secure'
})
</script>

<style lang="scss"></style>
