<template>
  <div id="tiles" class="mb-5">
    <div v-if="!loading" id="tile-data">
      <button
        class="button-link"
        :disabled="count <= perpage"
        @click="
          evt => {
            evt.preventDefault()
            changePage(false)
          }
        "
      >
        <no-ssr>
          <font-awesome-icon
            class="mr-2"
            style="max-width: 13px;"
            icon="arrow-left"
          />
        </no-ssr>
      </button>
      <button
        class="button-link"
        :disabled="count <= perpage"
        @click="
          evt => {
            evt.preventDefault()
            changePage(true)
          }
        "
      >
        <no-ssr>
          <font-awesome-icon
            class="mr-2"
            style="max-width: 13px;"
            icon="arrow-right"
          />
        </no-ssr>
      </button>
      <b-card-group deck>
        <no-ssr>
          <button
            v-for="(postval, index) in shownPosts"
            :key="`tile-${index}`"
            class="button-link"
            @click="
              evt => {
                evt.preventDefault()
                navigate(postval.id)
              }
            "
          >
            <b-card class="tile" no-body>
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
          </button>
        </no-ssr>
      </b-card-group>
    </div>
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
      imgUrl: cloudStorageURLs.posts,
      shownPosts: [],
      perpage: 2,
      loading: true
    }
  },
  async mounted() {
    /* eslint-disable */
    await this.updateCount()
    if (!this.currentIndex) {
      let currentindex = Math.floor((this.$store.state.tiles.perpage - this.perpage) / 2)
      if (this.count < this.perpage) {
        currentindex = 0
      }
      this.$store.commit('tiles/setIndex', {
        type: this.type,
        index: currentindex
      })
    }
    if (this.count !== 0 && this.count !== this.allPosts.length) {
      await this.initializePosts()
    }
    this.updateShownPosts()
  },
  computed: {
    currentIndex() {
      if (this.type === 'blog') {
        return this.$store.state.tiles.blogindex
      }
      return this.$store.state.tiles.projectindex
    },
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
    navigate(id) {
      // @ts-ignore
      if (process.client) {
        this.$router.push({
          path: `/${this.type}/${id}`
        })
      }
    },
    async updateShownPosts() {
      if (this.count === 0 || this.perpage < 1) {
        this.shownPosts = []
        this.loading = false
        return
      }
      const startpage = Math.floor(this.currentIndex / this.$store.state.tiles.perpage)
      const endpage = Math.ceil((this.currentIndex + this.perpage - 1) / this.$store.state.tiles.perpage)
      console.log(`start at ${this.currentIndex}`)
      console.log(`start page ${startpage}, end ${endpage}`)
      for (let i = startpage; i < endpage; i++) {
        if (!this.allPosts[i]) {
          await this.addPosts(i)
        }
      }
      const startpageindex = this.currentIndex % this.$store.state.tiles.perpage
      let endpageindex = (startpageindex + this.count) % this.$store.state.tiles.perpage
      if (this.count < this.$store.state.tiles.perpage) {
        endpageindex = this.count
      }
      console.log(`count ${this.count}, perpage: ${this.perpage}`)
      console.log(`start ${startpageindex} end page index ${endpageindex}`)
      const newShownPosts: any = []
      for (let i = startpage; i < endpage; i++) {
        console.log(`i ${i} ${endpage - 1}`)
        let start = i === startpage ? startpageindex : 0
        let end = i === (endpage - 1) ? endpageindex : this.$store.state.tiles.perpage
        console.log(`start ${start}, end ${end}`)
        console.log(this.allPosts[i])
        for (let j = start; j < end; j++) {
          console.log(`j ${j}`)
          const newPost: any = this.allPosts[i][j]
          newPost.title = decodeURIComponent(newPost.title)
          newPost.caption = decodeURIComponent(newPost.caption)
          newShownPosts.push(newPost)
          console.log(`j ${j}`)
        }
      }
      this.shownPosts = newShownPosts
      this.loading = false
    },
    changePage(increase) {
      console.log('change the page')
      let newindex = (increase ? this.currentIndex + 1 : this.currentIndex - 1) % this.count
      if (newindex < 0) newindex += this.count
      this.$store.commit('tiles/setIndex', {
        type: this.type,
        index: newindex
      })
      this.updateShownPosts()
    },
    updateCount() {
      return this.$store.dispatch('tiles/updateCount', {
        type: this.type
      }).then(res => {
        console.log(`got count ${res}`)
      }).catch(err => {
        console.log(err)
        this.$toasted.global.error({
          message: err
        })
      })
    },
    initializePosts() {
      return this.$store.dispatch('tiles/initializePosts', {
        type: this.type
      }).then(res => {
        console.log(`got init res ${res}`)
      }).catch(err => {
        console.log(err)
        this.$toasted.global.error({
          message: err
        })
      })
    },
    addPosts(page) {
      return this.$store.dispatch('tiles/addPosts', {
        type: this.type,
        page: page
      }).then(res => {
        console.log(`got add post res ${res}`)
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
  width: 200px;
  height: 200px;
}
.tile {
  text-align: center;
  max-width: 250px;
}
.zoom:hover {
  transform: scale(1.05);
}
</style>
