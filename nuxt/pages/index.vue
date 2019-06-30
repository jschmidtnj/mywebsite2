<template>
  <section class="container">
    <div>
      <logo />
      <h1>Under Construction</h1>
      <h2>Please check again later</h2>
    </div>
  </section>
</template>

<script lang="ts">
import Vue from 'vue'
import Logo from '~/components/Logo.vue'

export default Vue.extend({
  components: {
    Logo
  },
  data() {
    return {}
  },
  // @ts-ignore
  head() {
    return {
      title: 'Home'
    }
  },
  mounted() {
    /* eslint-disable */
    console.log('hello world!')
  }
})
</script>

<style lang="scss" scoped>
.text {
  color: white;
  font-size: 20px;
  position: absolute;
  top: 50%;
  left: 50%;
  -webkit-transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  text-align: center;
}

.container {
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.links {
  padding-top: 15px;
}
</style>
