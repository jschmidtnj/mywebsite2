<template>
  <b-card title="Blogs">
    <post-search type="blog" />
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
import PostSearch from '~/components/PostSearch.vue'
export default Vue.extend({
  name: 'BlogSearch',
  components: {
    PostSearch
  },
  data() {
    return {}
  },
  // @ts-ignore
  head() {
    return {
      title: 'Blogs'
    }
  }
})
</script>

<style lang="scss"></style>
