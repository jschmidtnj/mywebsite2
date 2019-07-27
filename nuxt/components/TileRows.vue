<template>
  <div id="tiles" class="mb-5">
    <b-card-group v-if="!loading" deck>
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
      loading: true
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
    navigate(id) {
      // @ts-ignore
      if (process.client) {
        this.$router.push({
          path: `/${this.type}/${id}`
        })
        window.location.reload(true)
      }
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
      this.shownPosts = this.allPosts
      this.loading = false
    },
    updateCount() {
      console.log(`got type ${this.type}`)
      return this.$store.dispatch('tiles/updateCount', {
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
      return this.$store.dispatch('tiles/initializePosts', {
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
      return this.$store.dispatch('tiles/addPosts', {
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
