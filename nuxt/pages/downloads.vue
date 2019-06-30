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

<script lang="ts">
import Vue from 'vue'
import { androidAppURL, cloudStorageURLs } from '~/assets/config'
export default Vue.extend({
  name: 'Downloads',
  data() {
    return {
      androidAppUrl: androidAppURL,
      windowsAppUrl: `${cloudStorageURLs.desktop}/win32`,
      macosAppUrl: `${cloudStorageURLs.desktop}/darwin`,
      linuxAppUrl: `${cloudStorageURLs.desktop}/linux`
    }
  },
  // @ts-ignore
  head() {
    return {
      title: 'Home'
    }
  }
})
</script>

<style lang="scss"></style>
