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
  </b-navbar>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'AdminNavbar'
})
</script>

<style lang="scss"></style>
