<template>
  <b-card title="Account Page">
    <p>user: {{ $auth.user }}</p>
    <b-btn @click="$auth.logout()">Logout</b-btn>
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  // @ts-ignore
  layout: 'secure',
  // @ts-ignore
  head() {
    return {
      title: 'Account'
    }
  }
})
</script>

<style lang="scss"></style>
