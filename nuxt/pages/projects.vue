<template>
  <post-search type="project" />
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
  },
  // @ts-ignore
  head() {
    return {
      title: 'Projects'
    }
  }
})
</script>

<style lang="scss"></style>
