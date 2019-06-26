<template>
  <post-search type="blog" />
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
  }
})
</script>

<style lang="scss">
</style>
