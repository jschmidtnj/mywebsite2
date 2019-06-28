<template>
  <post-data type="project" />
</template>

<script lang="ts">
import Vue from 'vue'
import PostData from '~/components/PostData.vue'
export default Vue.extend({
  name: 'Project',
  components: {
    PostData
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss"></style>
