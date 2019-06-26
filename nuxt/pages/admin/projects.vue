<template>
  <admin-post-edit type="project" />
</template>

<script lang="ts">
import Vue from 'vue'
import AdminPostEdit from '~/components/admin/PostEdit.vue'
export default Vue.extend({
  name: 'ProjectEdit',
  components: {
    AdminPostEdit
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss">
</style>
