<template>
  <div id="app">
    <atom-spinner :animation-duration="1000" :size="60" :color="'#ff1d5e'" />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { AtomSpinner } from 'epic-spinners'
export default Vue.extend({
  components: {
    AtomSpinner
  }
})
</script>
