<template>
  <b-navbar toggleable="lg" type="dark" variant="info">
    <b-navbar-brand href="/admin">Admin</b-navbar-brand>
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item href="/admin/blogs">Blogs</b-nav-item>
        <b-nav-item href="/admin/projects">Projects</b-nav-item>
      </b-navbar-nav>
    </b-collapse>
    <b-navbar-nav class="ml-auto">
      <b-nav-item-dropdown right>
        <template slot="button-content"
          ><em>User</em></template
        >
        <b-dropdown-item href="/account">Profile</b-dropdown-item>
        <b-dropdown-item href="#" @click="$auth.logout()"
          >Sign Out</b-dropdown-item
        >
      </b-nav-item-dropdown>
    </b-navbar-nav>
  </b-navbar>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'AdminNavbar'
})
</script>

<style lang="scss"></style>
