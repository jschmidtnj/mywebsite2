<template>
  <b-container class="mt-4">
    <b-card title="Blogs">
      <post-search type="blog" />
    </b-card>
  </b-container>
</template>

<script lang="js">
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
