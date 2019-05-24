<template>
  <section class="container">
    <div>
      <logo />
      <h1 class="title">joshua's website is awesome</h1>
      <h2 class="subtitle">personal website front-end nuxt</h2>
      <b-alert show>Example alert</b-alert>
      <no-ssr>
        <font-awesome-icon
          class="mr-2"
          style="max-width: 16px;"
          icon="search"
        />Example search icon
      </no-ssr>
      <p>{{ count }}</p>
      <div class="smallcontainer">
        <img
          src="https://mdbootstrap.com/img/Photos/Others/forest-sm.jpg"
          alt="Avatar"
          class="image"
        />
        <div class="overlay">
          <div class="text">Hello World</div>
        </div>
      </div>
      <div class="links">
        <a
          href="https://joshua-schmidt.github.io/joshuawebsite/"
          target="_blank"
          class="button--green"
          >Documentation</a
        >
        <a
          href="https://github.com/joshua-schmidt/joshuawebsite"
          target="_blank"
          class="button--grey"
          >GitHub</a
        >
      </div>
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
    return {
      count: 0
    }
  },
  mounted() {
    /* eslint-disable */
    console.log('hello world!')
    setInterval(() => {
      this.count++
    }, 1000)
  }
})
</script>

<style lang="scss" scoped>
.smallcontainer {
  position: relative;
  width: 50%;
}

.image {
  display: block;
  width: 100%;
  height: auto;
}

.overlay {
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100%;
  width: 100%;
  opacity: 0;
  transition: .5s ease;
  background-color: #008CBA;
}

.smallcontainer:hover .overlay {
  opacity: 1;
}

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

.title {
  font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
    'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>
