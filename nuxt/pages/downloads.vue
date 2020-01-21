<template>
  <b-card title="Downloads">
    <br />
    <b-list-group>
      <b-list-group-item>
        <b-link :href="androidAppUrl" class="card-link"
          >Download Android App</b-link
        >
      </b-list-group-item>
      <b-list-group-item>
        <b-link :href="windowsAppUrl" class="card-link"
          >Download Windows App</b-link
        >
      </b-list-group-item>
      <b-list-group-item>
        <b-link :href="macosAppUrl" class="card-link"
          >Download MacOS App</b-link
        >
      </b-list-group-item>
      <b-list-group-item>
        <b-link :href="linuxAppUrl" class="card-link"
          >Download Linux App</b-link
        >
      </b-list-group-item>
    </b-list-group>
  </b-card>
</template>

<script lang="js">
import Vue from 'vue'
import { androidAppURL, cloudStorageURLs } from '~/assets/config'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  name: 'Downloads',
  data() {
    return {
      androidAppUrl: androidAppURL,
      windowsAppUrl: `${cloudStorageURLs.desktop}/releases/win32`,
      macosAppUrl: `${cloudStorageURLs.desktop}/releases/darwin`,
      linuxAppUrl: `${cloudStorageURLs.desktop}/releases/linux`
    }
  },
  // @ts-ignore
  head() {
    const title = 'Downloads'
    const description =
      'download application for mac, linux, windows, or android'
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
