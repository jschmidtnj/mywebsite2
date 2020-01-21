<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="/">Joshua Short</b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item v-if="!loggedin" href="/signup">Signup</b-nav-item>
        <b-nav-item v-if="!loggedin" href="/login">Login</b-nav-item>
        <b-nav-item v-if="loggedin" href="/account">Account</b-nav-item>
      </b-navbar-nav>
      <b-navbar-nav v-if="loggedin" class="ml-auto">
        <b-nav-item-dropdown right>
          <b-dropdown-item @click="logout" href="#">
            Sign Out
          </b-dropdown-item>
        </b-nav-item-dropdown>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script lang="js">
import Vue from 'vue'
export default Vue.extend({
  name: 'Navbar',
  computed: {
    loggedin() {
      return this.$store.state.auth && this.$store.state.auth.loggedIn
    }
  },
  methods: {
    logout(evt) {
      /* eslint-disable */
      evt.preventDefault()
      this.$store.commit('auth/logout')
      console.log(`layout name ${this.$nuxt.$data.layoutName}`)
      if (this.$nuxt.$data.layoutName === 'secure') {
        this.$router.push({
          path: '/'
        })
      }
    }
  }
})
</script>

<style lang="scss"></style>
