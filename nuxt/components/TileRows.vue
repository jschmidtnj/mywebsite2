<template>
  <div id="tiles" class="mb-5">
    <div v-if="!loading">
      <b-row v-for="(rowval, rowindex) in shownPosts" :key="`row-${rowindex}`">
        <b-col
          v-for="(postval, colindex) in shownPosts[rowindex]"
          :key="`post-${rowindex}-${colindex}`"
          sm
          style="padding:0px;"
        >
          <a v-if="postval" :href="`/${type}/${postval.id}`">
            <b-card class="tile rounded-0" no-body>
              <b-card-body
                class="tile-body zoom"
                @mouseenter="
                  selected[rowindex * shownPosts[0].length + colindex] = true
                  $forceUpdate()
                "
                @mouseleave="
                  selected[rowindex * shownPosts[0].length + colindex] = false
                  $forceUpdate()
                "
              >
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
                  class="tile-img rounded-0"
                />
                <b-container
                  v-if="selected[rowindex * shownPosts[0].length + colindex]"
                  :style="{ 'background-color': getColor(rowindex, colindex) }"
                  class="main-overlay"
                >
                  <div class="text-overlay">
                    <b-card-title class="white-color">{{
                      postval.title
                    }}</b-card-title>
                    <b-card-sub-title class="white-color">{{
                      postval.caption
                    }}</b-card-sub-title>
                  </div>
                </b-container>
              </b-card-body>
            </b-card>
          </a>
          <div v-else></div>
        </b-col>
      </b-row>
      <p v-if="shownPosts.length === 0">No {{ `${type}s` }} found.</p>
    </div>

    <loading v-else />
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Loading from '~/components/ComponentLoading.vue'
import { validTypes, cloudStorageURLs } from '~/assets/config'
const defaultOpacity = 20 // %
const numPerRow = 2
const defaultOpacityHex = Math.round((defaultOpacity / 100.0) * 255)
  .toString(16)
  .toUpperCase()
export default Vue.extend({
  name: 'TileRows',
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
      imgUrl: cloudStorageURLs.posts,
      shownPosts: [],
      loading: true,
      selected: []
    }
  },
  async mounted() {
    /* eslint-disable */
    await this.updateCount()
    if (this.count !== 0 && this.count !== this.allPosts.length) {
      await this.initializePosts()
    }
    this.updateShownPosts()
  },
  computed: {
    count() {
      if (this.type === 'blog') {
        return this.$store.state.tiles.blogcount
      }
      return this.$store.state.tiles.projectcount
    },
    allPosts() {
      if (this.type === 'blog') {
        return this.$store.state.tiles.blogs
      }
      return this.$store.state.tiles.projects
    }
  },
  methods: {
    getColor(rowindex, colindex) {
      if (this.selected[rowindex * this.shownPosts[0].length + colindex]) {
        let colorHex = this.shownPosts[rowindex][colindex].color
        if (colorHex.length === 7) {
          colorHex = colorHex.concat(defaultOpacityHex)
        }
        return colorHex
      }
      return 'white'
    },
    async updateShownPosts() {
      if (this.count === 0) {
        this.shownPosts = []
        this.loading = false
        return
      }
      const endpage = Math.ceil(this.count / this.$store.state.tiles.perpage)
      for (let i = 0; i < endpage; i++) {
        if (!this.allPosts[i]) {
          await this.addPosts(i)
        }
      }
      let newShownPosts: any = []
      newShownPosts.push([])
      let currentIndex = 0
      for (let i = 0; i < this.allPosts.length; i++) {
        for (let j = 0; j < this.allPosts[i].length; j++) {
          const newPost: any = this.allPosts[i][j]
          newPost.title = decodeURIComponent(newPost.title)
          newPost.caption = decodeURIComponent(newPost.caption)
          newPost.color = decodeURIComponent(newPost.color)
          if (newShownPosts[currentIndex].length === numPerRow) {
            newShownPosts.push([newPost])
            currentIndex++
          } else {
            newShownPosts[currentIndex].push(newPost)
          }
        }
      }
      for (let i = newShownPosts[currentIndex].length; i < numPerRow; i++) {
        newShownPosts[currentIndex].push(null)
      }
      this.shownPosts = newShownPosts
      this.loading = false
    },
    updateCount() {
      return this.$store
        .dispatch('tiles/updateCount', {
          type: this.type
        })
        .then(res => {
          console.log(`got res ${res}`)
          this.selected = new Array(this.count).fill(false)
        })
        .catch(err => {
          console.log(err)
          this.$toasted.global.error({
            message: err
          })
        })
    },
    initializePosts() {
      return this.$store
        .dispatch('tiles/initializePosts', {
          type: this.type
        })
        .then(res => {
          console.log(`got res ${res}`)
        })
        .catch(err => {
          console.log(err)
          this.$toasted.global.error({
            message: err
          })
        })
    },
    addPosts(page) {
      return this.$store
        .dispatch('tiles/addPosts', {
          type: this.type,
          page: page
        })
        .then(res => {
          console.log(`got res ${res}`)
        })
        .catch(err => {
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
.white-color {
  color: white;
}
.tile-img img {
  object-fit: cover;
  width: 100%;
  height: 200px;
  position: relative;
}
.tile {
  overflow: hidden;
}
.tile-body {
  text-align: center;
  width: 100%;
  height: 200px;
  padding: 0;
}
.main-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height:100%;
  z-index: 9999;
}
.text-overlay {
  margin-top: 10%;
  width: 100%;
  height: 100%;
}
.zoom:hover {
  transform: scale(1.05);
  -moz-transform: scale(1.05);
  -webkit-transform: scale(1.05);
}
</style>
