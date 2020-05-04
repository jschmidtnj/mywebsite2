<template>
  <b-container class="mt-4">
    <b-card title="Projects">
      <post-search type="project" />
    </b-card>
  </b-container>
</template>

<script lang="js">
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
