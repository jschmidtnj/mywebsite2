<template>
  <section class="container">
    <div>
      <b-jumbotron
        header="Joshua Schmidt"
        lead="Programmer, Computer Engineer, Developer"
        class="mt-4"
      >
        <b-button variant="primary" href="/about">About</b-button>
      </b-jumbotron>
      <h3>Blogs</h3>
      <tile-rows type="blog" />
      <h3>Projects</h3>
      <tile-rows type="project" />
    </div>
  </section>
</template>

<script lang="ts">
import Vue from 'vue'
import TileRows from '~/components/TileRows.vue'

export default Vue.extend({
  components: {
    TileRows
  },
  data() {
    return {}
  },
  // @ts-ignore
  head() {
    return {
      title: 'Home'
    }
  },
  mounted() {}
})
</script>

<style lang="scss"></style>
