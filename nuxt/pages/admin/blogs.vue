<template>
  <admin-post-edit type="blog" />
</template>

<script lang="ts">
import Vue from 'vue'
import AdminPostEdit from '~/components/admin/PostEdit.vue'
export default Vue.extend({
  name: 'BlogEdit',
  // @ts-ignore
  layout: 'admin',
  components: {
    AdminPostEdit
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss"></style>
