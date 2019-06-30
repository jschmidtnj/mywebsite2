<template>
  <post-data type="blog" />
</template>

<script lang="ts">
import Vue from 'vue'
import PostData from '~/components/PostData.vue'
export default Vue.extend({
  name: 'Blog',
  components: {
    PostData
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss"></style>
