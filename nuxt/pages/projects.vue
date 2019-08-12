<template>
  <b-card title="Projects">
    <post-search type="project" />
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
import PostSearch from '~/components/PostSearch.vue'
export default Vue.extend({
  name: 'ProjectSearch',
  components: {
    PostSearch
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss"></style>
