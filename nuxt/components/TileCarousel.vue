<template>
  <div :id="`tile-carousel-${type}`">
    <div v-if="!loading" id="tile-data">
      <div v-if="count > perpage" id="navigation-buttons" class="mb-3">
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
              class="mr-2 arrow-size-carousel"
              icon="chevron-left"
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
              class="mr-2 arrow-size-carousel"
              icon="chevron-right"
            />
          </no-ssr>
        </button>
      </div>
      <b-card-group
        ref="carouselContent"
        deck
        class="scrolling-wrapper flex-row flex-nowrap"
      >
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
            <b-card class="tile m-2" no-body>
              <b-card-body class="tile zoom p-0">
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
                  <b-card-title title-tag="h5">
                    {{ postval.caption }}
                  </b-card-title>
                </b-container>
              </b-card-body>
            </b-card>
          </button>
        </no-ssr>
      </b-card-group>
      <div v-if="window.width < count * 200" id="scroll-buttons" class="mt-3">
        <button class="button-link" @click="swipeLeft">
          <no-ssr>
            <font-awesome-icon
              class="mr-2 arrow-size-carousel"
              icon="chevron-left"
            />
          </no-ssr>
        </button>
        <button class="button-link" @click="swipeRight">
          <no-ssr>
            <font-awesome-icon
              class="mr-2 arrow-size-carousel"
              icon="chevron-right"
            />
          </no-ssr>
        </button>
      </div>
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
      perpage: 8,
      loading: true,
      window: {
        width: 0,
        height: 0
      }
    }
  },
  destroyed() {
    window.removeEventListener('resize', this.handleResize)
  },
  async mounted() {
    /* eslint-disable */
    window.addEventListener('resize', this.handleResize)
    this.handleResize()
    await this.updateCount()
    if (!this.currentIndex) {
      let currentindex = 0
      if (this.count > this.perpage) {
        let currentindex = Math.floor((this.$store.state.tiles.perpage - this.perpage) / 2)
        if (currentindex < 0) currentindex = this.count - 1
        else if (currentindex >= this.count) currentindex = 0
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
    handleResize() {
      this.window.width = window.innerWidth;
      this.window.height = window.innerHeight;
    },
    /**
     * scrollTo - Horizontal Scrolling
     * @param {(HTMLElement ref)} element - Scroll Container
     * @param {number} scrollPixels - pixel to scroll
     * @param {number} duration -  Duration of scrolling animation in millisec
     */
    scrollTo(element, scrollPixels, duration) {
      // Get the start timestamp
      const startTime =
        "now" in window.performance
          ? performance.now()
          : new Date().getTime();
      const scroll = (timestamp) => {
        //Calculate the timeelapsed
        const timeElapsed = timestamp - startTime;
        //Calculate progress 
        const progress = Math.min(timeElapsed / duration, 1);
        //Set the scrolleft
        element.scrollLeft = scrollPos + scrollPixels * progress;
        //Check if elapsed time is less then duration then call the requestAnimation, otherwise exit
        if (timeElapsed < duration) {
          //Request for animation
          window.requestAnimationFrame(scroll);
        } else {
          return;
        }
      }
      const scrollPos = element.scrollLeft;
      // Condition to check if scrolling is required
      if ( !( (scrollPos === 0 || scrollPixels > 0) && (element.clientWidth + scrollPos === element.scrollWidth || scrollPixels < 0))) 
      {
        //Call requestAnimationFrame on scroll function first time
        window.requestAnimationFrame(scroll);
      }
    },
    swipeLeft() {
      const content = this.$refs.carouselContent;
      this.scrollTo(content, -300, 500);
    },
    swipeRight() {
      const content = this.$refs.carouselContent;
      this.scrollTo(content, 300, 500);
    },
    navigate(id) {
      // @ts-ignore
      if (process.client) {
        this.$router.push({
          path: `/${this.type}/${id}`
        })
      }
    },
    async updateShownPosts() {
      if (this.count === 0 || this.perpage < 1 || this.currentIndex >= this.count) {
        this.shownPosts = []
        this.loading = false
        return
      }
      const startpage = Math.floor(this.currentIndex / this.$store.state.tiles.perpage)
      const endpage = Math.ceil((this.currentIndex + this.perpage) / this.$store.state.tiles.perpage)
      const startpageindex = this.currentIndex % this.$store.state.tiles.perpage
      const allPostsLen = Math.ceil(this.count / this.$store.state.tiles.perpage)
      const allPostsIndexLen = this.count < this.$store.state.tiles.perpage ? this.count : this.$store.state.tiles.perpage
      const shownPostsLen = this.count < this.perpage ? this.count : this.perpage
      const newShownPosts: any = []
      for (let i = startpage; i < endpage; i++) {
        let start = i === startpage ? startpageindex : 0
        await this.addPosts(i % allPostsLen)
        for (let j = start; (j < allPostsIndexLen || (i === endpage - 1 && j < allPostsIndexLen * 2)) && newShownPosts.length < shownPostsLen; j++) {
          const newPost: any = this.allPosts[i % allPostsLen][j % allPostsIndexLen]
          Object.keys(newPost).forEach(key => {
            if (newPost[key] instanceof String)
              newPost[key] = decodeURIComponent(newPost[key]);
          })
          newShownPosts.push(newPost)
        }
      }
      this.shownPosts = newShownPosts
      this.loading = false
    },
    changePage(increase) {
      let newindex = (increase ? this.currentIndex + 1 : this.currentIndex - 1)
      if (newindex < 0) newindex = this.count - 1
      else if (newindex >= this.count) newindex = 0
      console.log(`change the page to ${newindex}`)
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
.arrow-size-carousel {
  font-size: 1.5rem;
}
.tile-img {
  object-fit: cover;
  width: 150px !important;
  height: 150px !important;
}
.tile {
  text-align: center;
  max-width: 250px;
  min-width: 150px;
}
.zoom:hover {
  transform: scale(1.05);
}
.scrolling-wrapper {
  overflow-x: scroll;
  overflow-y: visible;
  white-space: nowrap;
  width: 100%;
  -webkit-overflow-scrolling: touch;
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
