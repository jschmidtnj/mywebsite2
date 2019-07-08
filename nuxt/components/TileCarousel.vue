<template>
  <div id="tiles">
    <b-container v-if="posts.length > 0">
      <b-card
        v-for="(postval, index) in posts"
        :key="`tile-${index}`"
        :title="postval.title"
        :img-src="postval.tile"
        :img-alt="postval.title"
        img-top
      >
        <b-card-text>
          {{ postval.caption }}
        </b-card-text>
      </b-card>
    </b-container>
    <atom-spinner
      v-else
      class="centered"
      :animation-duration="1000"
      :size="60"
      :color="'#ff1d5e'"
    />
    <!-- need arrow right and left from fontawesome, and start in the middle
    then if you move too far to one side, query for more posts and save the
    current posts in case you cycle through again -->
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { AtomSpinner } from 'epic-spinners/dist/lib/epic-spinners.min.js'
import { validTypes } from '~/assets/config'
export default Vue.extend({
  name: 'TileCarousel',
  props: {
    type: {
      type: String,
      default: null,
      required: true,
      validator: val => validTypes.includes(String(val))
    }
  },
  components: {
    AtomSpinner
  },
  mounted() {
    /* eslint-disable */
    if (this.posts.length === 0) {
      this.updatePosts()
    }
  },
  computed: {
    posts() {
      if (this.type === 'blog') {
        return this.$store.state.tilecarousel.blogs
      }
      return this.$store.state.tilecarousel.projects
    }
  },
  methods: {
    updatePosts() {
      this.$store.dispatch('tilecarousel/updateCarousel').then(res => {
        console.log(`got res ${res}`)
      }).catch(err => {
        this.$toasted.global.error({
          message: err
        })
      })
    }
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
