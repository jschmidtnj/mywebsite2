<template>
  <div id="tiles" class="mb-5">
    <b-card-group v-if="posts.length > 0" deck>
      <no-ssr>
        <b-card
          v-for="(postval, index) in posts"
          :key="`tile-${index}`"
          class="tile"
          no-body
        >
          <b-card-body class="tile zoom">
            <b-card-img-lazy
              :src="
                `${imgUrl}/${
                  type === 'blog' ? 'blogimages' : 'projectimages'
                }/${encodeURI(postval.tileimage)}/original`
              "
              :blank-src="
                `${imgUrl}/${
                  type === 'blog' ? 'blogimages' : 'projectimages'
                }/${encodeURI(postval.tileimage)}/blur`
              "
              :alt="postval.title"
              class="tile-img"
            />
            <b-container>
              <b-card-title>
                {{ postval.title }}
              </b-card-title>
              <b-card-sub-title>
                {{ postval.caption }}
              </b-card-sub-title>
            </b-container>
          </b-card-body>
        </b-card>
      </no-ssr>
    </b-card-group>
    <loading v-else />
    <!-- need arrow right and left from fontawesome, and start in the middle
    then if you move too far to one side, query for more posts and save the
    current posts in case you cycle through again -->
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Loading from '~/components/ComponentLoading.vue'
import { validTypes, cloudStorageURLs } from '~/assets/config'
export default Vue.extend({
  name: 'TileCarousel',
  components: {
    Loading
  },
  props: {
    type: {
      type: String,
      default: null,
      required: true,
      validator: val => validTypes.includes(String(val))
    }
  },
  data() {
    return {
      imgUrl: cloudStorageURLs.posts
    }
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
      this.$store.dispatch('tilecarousel/updateCarousel', {
        type: this.type
      }).then(res => {
        console.log(`got res ${res}`)
      }).catch(err => {
        console.log(err)
        this.$toasted.global.error({
          message: err
        })
      })
    }
  }
})
</script>

<style lang="scss">
.tile-img {
  object-fit: cover;
  width: 300px;
  height: 300px;
}
.tile {
  text-align: center;
  max-width: 350px;
}
.zoom:hover {
  transform: scale(1.05);
}
</style>
