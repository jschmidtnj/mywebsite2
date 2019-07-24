<template>
  <div>
    <p>Admin Page</p>
    <b-btn block @click="logout">Logout</b-btn>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  // @ts-ignore
  layout: 'admin',
  methods: {
    logout(evt) {
      evt.preventDefault()
      this.$store.commit('auth/logout')
      this.$router.push({
        path: '/login'
      })
    }
  }
})
</script>

<style lang="scss"></style>
