<template>
  <section class="container">
    <div>
      <h1>Under Construction</h1>
      <h2>Please check again later</h2>
      <hr />
      <tile-rows type="blog" />
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
