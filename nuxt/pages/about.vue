<template>
  <b-container>
    <b-jumbotron header="About me" class="mt-4"></b-jumbotron>
    <b-list-group>
      <b-list-group-item
        >Student at Stevens Institute of Technology</b-list-group-item
      >
      <b-list-group-item>Ultimate Frisbee Player</b-list-group-item>
      <b-list-group-item>Computer Engineering Major</b-list-group-item>
      <b-list-group-item>Electronics guru</b-list-group-item>
      <b-list-group-item>Experienced Software Developer</b-list-group-item>
    </b-list-group>
  </b-container>
</template>

<script lang="ts">
import Vue from 'vue'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  name: 'About',
  // @ts-ignore
  head() {
    const title = 'About'
    const description = 'description of this website'
    const image = `${seo.url}/icon.png`
    return {
      title: title,
      meta: [
        { property: 'og:title', content: title },
        { property: 'og:description', content: description },
        {
          property: 'og:image',
          content: image
        },
        { name: 'twitter:title', content: title },
        {
          name: 'twitter:description',
          content: description
        },
        {
          name: 'twitter:image',
          content: image
        },
        { hid: 'description', name: 'description', content: description }
      ]
    }
  }
})
</script>

<style lang="scss"></style>
