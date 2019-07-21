<template>
  <atom-spinner
    class="centered"
    :animation-duration="1000"
    :size="60"
    :color="'#ff1d5e'"
  />
</template>

<script lang="ts">
import Vue from 'vue'
import { AtomSpinner } from 'epic-spinners/dist/lib/epic-spinners.min.js'
export default Vue.extend({
  components: {
    AtomSpinner
  }
})
</script>

<style lang="scss">
@import '~/node_modules/epic-spinners/dist/lib/epic-spinners.min.css';

.centered {
  display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>
