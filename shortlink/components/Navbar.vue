<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="/">Joshua Short</b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item v-if="!loggedin" :href="`${mainurl}/signup`"
          >Signup</b-nav-item
        >
        <b-nav-item
          v-if="!loggedin"
          :href="`${mainurl}/login?redirect_uri=${redirecturl}`"
          >Login</b-nav-item
        >
      </b-navbar-nav>
      <b-navbar-nav v-if="loggedin" class="ml-auto">
        <b-nav-item-dropdown right>
          <b-dropdown-item href="#" @click="$auth.logout()">
            Sign Out
          </b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script lang="ts">
import Vue from 'vue'
const mainurl = process.env.mainurl ? process.env.mainurl : ''
const currenturl = process.env.seoconfig
  ? JSON.parse(process.env.seoconfig).url
  : ''
const redirecturl = encodeURIComponent(`${currenturl}/callback`)
export default Vue.extend({
  name: 'Navbar',
  data() {
    return {
      mainurl: mainurl,
      redirecturl: redirecturl
    }
  },
  computed: {
    loggedin() {
      return this.$store.state.auth.loggedIn
    }
  }
})
</script>

<style lang="scss"></style>
