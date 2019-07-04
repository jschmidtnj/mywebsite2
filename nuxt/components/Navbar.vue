<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="/">Joshua Schmidt</b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item href="/blogs">Blogs</b-nav-item>
        <b-nav-item href="/projects">Projects</b-nav-item>
        <b-nav-item href="/downloads">Downloads</b-nav-item>
        <b-nav-item v-if="!loggedin" href="/signup">Signup</b-nav-item>
        <b-nav-item v-if="!loggedin" href="/login">Login</b-nav-item>
      </b-navbar-nav>
      <b-navbar-nav v-if="loggedin" class="ml-auto">
        <b-nav-item-dropdown right>
          <template slot="button-content">
            <em>User</em>
          </template>
          <b-dropdown-item href="/account">Profile</b-dropdown-item>
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
export default Vue.extend({
  name: 'Navbar',
  data() {
    return {}
  },
  computed: {
    loggedin() {
      return this.$auth.user
    }
  }
})
</script>

<style lang="scss"></style>
