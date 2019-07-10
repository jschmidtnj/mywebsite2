<template>
  <div id="tiles" class="mb-5">
    <b-button
      variant="primary"
      @click="
        evt => {
          evt.preventDefault()
          changePage(true)
        }
      "
      >Forward</b-button
    >
    <b-button
      variant="primary"
      @click="
        evt => {
          evt.preventDefault()
          changePage(false)
        }
      "
      >Backwards</b-button
    >
    <b-card-group v-if="shownPosts.length > 0" deck>
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
      perpage: 4
    }
  },
  async mounted() {
    /* eslint-disable */
    if (!this.currentIndex) {
      this.$store.commit('setIndex', {
        type: this.type,
        index: Math.floor((
          this.$store.state.tilecarousel.perpage - this.perpage) / 2
        )
      })
    }
    await this.updateCount()
    if (this.count !== 0 && this.count !== this.allPosts.length) {
      await this.initializePosts()
    }
    this.updateShownPosts()
  },
  computed: {
    currentIndex() {
      if (this.type === 'blog') {
        return this.$store.state.tilecarousel.blogindex
      }
      return this.$store.state.tilecarousel.projectindex
    },
    count() {
      if (this.type === 'blog') {
        return this.$store.state.tilecarousel.blogcount
      }
      return this.$store.state.tilecarousel.projectcount
    },
    allPosts() {
      if (this.type === 'blog') {
        return this.$store.state.tilecarousel.blogs
      }
      return this.$store.state.tilecarousel.projects
    }
  },
  methods: {
    navigate(id) {
      // @ts-ignore
      if (process.client) {
        this.$router.push({
          path: `/${this.type}?id=${id}`
        })
        window.location.reload(true)
      }
    },
    async updateShownPosts() {
      if (this.count === 0) {
        this.shownPosts = []
        return
      }
      const startpage = Math.floor(this.currentIndex / this.$store.state.tilecarousel.perpage)
      const endpage = Math.ceil((this.currentIndex + this.perpage) / this.$store.state.tilecarousel.perpage)
      for (let i = startpage; i < endpage; i++) {
        if (!this.allPosts[i]) {
          await this.addPosts(i)
        }
      }
      const startpageindex = this.currentIndex % (this.$store.state.tilecarousel.perpage - 1)
      const endpageindex = (startpageindex + this.perpage) % (this.$store.state.tilecarousel.perpage - 1)
      const newShownPosts: any = []
      for (let i = startpage; i < endpage; i++) {
        let start = i === startpage ? startpageindex : 0
        let end = i === endpage ? (this.$store.state.tilecarousel.perpage - 1) : endpageindex
        for (let j = start; j <= end; j++) {
          newShownPosts.push(this.allPosts[i][j])
        }
      }
      this.shownPosts = newShownPosts
    },
    changePage(increase) {
      let newindex = (increase ? this.currentIndex + 1 : this.currentIndex - 1) % this.count
      if (newindex < 0) newindex += this.count
      this.$store.commit('setIndex', {
        type: this.type,
        index: newindex
      })
      this.updateShownPosts()
    },
    updateCount() {
      console.log(`got type ${this.type}`)
      return this.$store.dispatch('tilecarousel/updateCount', {
        type: this.type
      }).then(res => {
        console.log(`got res ${res}`)
      }).catch(err => {
        console.log(err)
        this.$toasted.global.error({
          message: err
        })
      })
    },
    initializePosts() {
      return this.$store.dispatch('tilecarousel/initializePosts', {
        type: this.type
      }).then(res => {
        console.log(`got res ${res}`)
      }).catch(err => {
        console.log(err)
        this.$toasted.global.error({
          message: err
        })
      })
    },
    addPosts(page) {
      return this.$store.dispatch('tilecarousel/addPosts', {
        type: this.type,
        page: page
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
.button-link {
  display: inline-block;
  position: relative;
  background-color: transparent;
  cursor: pointer;
  border: 0;
  padding: 0;
  font: inherit;
}
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
