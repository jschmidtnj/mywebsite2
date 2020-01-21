<template>
  <div>
    <div id="post-data" v-if="post">
      <div id="post-content">
        <b-container class="hero-body">
          <b-row>
            <b-col>
              <b-img-lazy
                v-if="post.heroimage"
                :blank-src="
                  `${postCdn}/${
                    type === 'blog'
                      ? staticstorageindexes.blogfiles
                      : staticstorageindexes.projectfiles
                  }/${post.id}/${post.heroimage.id + paths.blur}`
                "
                :src="
                  `${postCdn}/${
                    type === 'blog'
                      ? staticstorageindexes.blogfiles
                      : staticstorageindexes.projectfiles
                  }/${post.id}/${post.heroimage.id + paths.original}`
                "
                :alt="post.heroimage.name"
                class="hero-img m-0"
              ></b-img-lazy>
              <div class="main-overlay">
                <div class="text-overlay">
                  <!-- add text overlay here -->
                </div>
              </div>
            </b-col>
          </b-row>
        </b-container>
        <b-container id="header-container" v-if="post">
          <h1>{{ post.title }}</h1>
          <p>{{ post.author }}</p>
          <p v-if="post.id">
            {{ formatDate(mongoidToDate(post.id), 'M/D/YYYY') }}
          </p>
          <p>{{ post.views }}</p>
          <a :href="`${shortlinkurl}/${post.shortlink}`">
            {{ `${shortlinkurl}/${post.shortlink}` }}
          </a>
          <p class="orange-text">
            {{
              post.categories
                .map(category => decodeURIComponent(category))
                .join(' | ')
            }}
          </p>
          <hr />
        </b-container>
        <b-container id="content-container" v-if="post">
          <vue-markdown
            :source="post.content"
            @rendered="updateMarkdown"
            class="markdown"
          />
        </b-container>
      </div>
      <b-container>
        <tile-carousel :type="type" class="mt-5 mb-2" />
      </b-container>
    </div>
    <loading v-else />
  </div>
</template>

<script lang="js">
import Vue from 'vue'
import { format } from 'date-fns'
import VueMarkdown from 'vue-markdown'
import Prism from 'prismjs'
import LazyLoad from 'vanilla-lazyload'
import Loading from '~/components/PageLoading.vue'
import TileCarousel from '~/components/TileCarousel.vue'
import {
  validTypes,
  cloudStorageURLs,
  staticstorageindexes,
  paths
} from '~/assets/config'
const lazyLoadInstance = new LazyLoad({
  elements_selector: '.lazy'
})
// @ts-ignore
const ampurl = process.env.ampurl
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
// @ts-ignore
const shortlinkurl = process.env.shortlinkurl
export default Vue.extend({
  name: 'Post',
  components: {
    VueMarkdown,
    Loading,
    TileCarousel
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
      id: null,
      post: null,
      shortlinkurl: shortlinkurl,
      postCdn: cloudStorageURLs.posts,
      staticstorageindexes: staticstorageindexes,
      paths: paths
    }
  },
  /* eslint-disable */
  mounted() {
    if (this.$route.params && this.$route.params.id) {
      this.id = this.$route.params.id
      this.$axios
        .get('/graphql', {
          params: {
            query: `{post(type:"${encodeURIComponent(
              this.type
            )}",id:"${encodeURIComponent(this.id)}",cache:${(!(
              this.$store.state.auth.user &&
              this.$store.state.auth.user.type === 'admin'
            )).toString()}){title caption content id author views shortlink heroimage{name id} tileimage{id} categories tags}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const post = res.data.data.post
                Object.keys(post).forEach(key => {
                  if (typeof post[key] === 'string')
                    post[key] = decodeURIComponent(post[key])
                })
                this.post = post
                // update title for spa
                document.title = this.post.title
              } else if (res.data.errors) {
                this.$nuxt.error({
                  statusCode: 404,
                  message: `found errors: ${JSON.stringify(res.data.errors)}`
                })
              } else {
                this.$nuxt.error({
                  statusCode: 404,
                  message: 'could not find data or errors'
                })
              }
            } else {
              this.$nuxt.error({
                statusCode: 404,
                message: 'could not get data'
              })
            }
          } else {
            this.$nuxt.error({
              statusCode: 404,
              message: `status code of ${res.status}`
            })
          }
        })
        .catch(err => {
          this.$nuxt.error({
            statusCode: 404,
            message: `got error: ${err}`
          })
        })
    } else {
      this.$nuxt.error({
        statusCode: 404,
        message: 'could not find id in params'
      })
    }
  },
  // @ts-ignore
  head() {
    const title = this.post ? this.post.title : this.type
    const description = this.post ? this.post.caption : this.type
    const meta = [
      { property: 'og:title', content: title },
      { property: 'og:description', content: description },
      { name: 'twitter:title', content: title },
      {
        name: 'twitter:description',
        content: description
      },
      { hid: 'description', name: 'description', content: description }
    ]
    const script = []
    if (this.post) {
      const image = `${cloudStorageURLs.posts}/${
        this.type === 'blog'
          ? this.staticstorageindexes.blogfiles
          : this.staticstorageindexes.projectfiles
      }/${this.post.id}/${this.post.tileimage.id + this.paths.original}`
      meta.push({
        property: 'og:image',
        content: image
      })
      meta.push({
        name: 'twitter:image',
        content: image
      })
      const date = this.formatDate(
        this.mongoidToDate(this.post.id),
        'YYYY-MM-DD'
      )
      script.push({
        innerHTML: JSON.stringify({
          '@context': 'https://schema.org',
          '@type': 'BlogPosting',
          headline: this.post.title,
          alternativeHeadline: this.post.caption,
          image: image,
          editor: this.post.author,
          genre: this.post.categories.join(' '),
          keywords: this.post.tags.join(' '),
          wordcount: this.post.content.length,
          publisher: seo.url,
          url: seo.url,
          datePublished: date,
          dateCreated: date,
          dateModified: date,
          description: this.post.caption,
          articleBody: this.post.content,
          author: {
            '@type': 'Person',
            name: this.post.author
          }
        }),
        type: 'application/ld+json'
      })
    }
    return {
      title: title,
      meta: meta,
      link: [
        {
          rel: 'amphtml',
          href: `${ampurl}/blog/${this.$route.query.id}`
        }
      ],
      __dangerouslyDisableSanitizers: ['script'],
      script: script
    }
  },
  methods: {
    updateMarkdown() {
      this.$nextTick(() => {
        Prism.highlightAll()
        if (lazyLoadInstance) {
          console.log('update lazyload')
          lazyLoadInstance.update()
        }
      })
    },
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    },
    mongoidToDate(id) {
      return parseInt(id.substring(0, 8), 16) * 1000
    }
  }
})
</script>

<style lang="scss">
#content-container {
  padding-left: 0;
  padding-right: 0;
}
#content-container p,
h1,
h2,
h3,
h4,
h5,
h6 {
  padding-right: 15px;
  padding-left: 15px;
}
#post-data {
  display: flex;
  min-height: 90vh;
  flex-direction: column;
}
#post-content {
  flex: 1;
}
@media (min-width: 1200px) {
  .container {
    max-width: 1400px;
  }
}
.white-color {
  color: white;
}
.hero-img {
  object-fit: cover;
  width: 100%;
  // set max height for image
  max-height: 40em;
  position: relative;
}
.hero-body {
  overflow: hidden;
  text-align: center;
  width: 100%;
  // set max height for image
  max-height: 40em;
  padding: 0;
}
.main-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 9999;
  // add gradiant to show text clearly
  // background: linear-gradient(rgba(0, 0, 0, 0.2), rgba(0, 0, 0, 0.2));
}
.text-overlay {
  padding-top: 10%;
  height: 100%;
}
</style>
