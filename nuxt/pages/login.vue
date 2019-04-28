<template>
  <div>
    <p>login</p>
    <b-btn block @click="login">Login with Google</b-btn>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'Login',
  // @ts-ignore
  middleware: 'loginredirect',
  methods: {
    login() {
      /* eslint-disable */
      this.$auth.loginWith('google').catch(e => {
        console.log(e)
      })
    }
  }
})
</script>

<style lang="scss"></style>
