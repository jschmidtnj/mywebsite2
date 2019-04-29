<template>
  <div>
    <p>login</p>
    <b-btn block @click="logingoogle">Login with Google</b-btn>
    <b-btn block @click="loginlocal">Login with local</b-btn>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'Login',
  // @ts-ignore
  middleware: 'loginredirect',
  methods: {
    logingoogle() {
      /* eslint-disable */
      this.$auth.loginWith('google').catch(e => {
        console.log(e)
      })
    },
    loginlocal() {
      this.$auth.loginWith('local', {
        data: {
          username: 'test',
          password: 'test'
        }
      })
    }
  }
})
</script>

<style lang="scss"></style>
