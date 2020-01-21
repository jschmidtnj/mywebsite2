<template>
  <atom-spinner
    :animation-duration="1000"
    :size="60"
    :color="'#ff1d5e'"
    class="centered"
  />
</template>

<script lang="js">
import Vue from 'vue'
import { AtomSpinner } from 'epic-spinners/dist/lib/epic-spinners.min.js'
export default Vue.extend({
  components: {
    AtomSpinner
  }
})
</script>

<style lang="scss">
.centered {
  display: block;
  margin-left: auto;
  margin-right: auto;
}
</style>
