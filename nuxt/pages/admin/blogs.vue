<template>
  <admin-post-edit type="blog" />
</template>

<script lang="ts">
import Vue from 'vue'
import AdminPostEdit from '~/components/admin/PostEdit.vue'
export default Vue.extend({
  name: 'BlogEdit',
  components: {
    AdminPostEdit
  },
  data() {
    return {}
  }
})
</script>

<style lang="scss"></style>
