<template>
  <b-container class="main-container" fluid>
    <b-row>
      <b-col md="6" class="my-1">
        <b-form-group label-cols-sm="3" label="search" class="mb-0">
          <b-input-group>
            <b-form-input
              v-model="search"
              placeholder="Type to Search"
              @keyup.enter.native="
                (evt) => {
                  evt.preventDefault()
                  currentPage = 1
                  searchPosts()
                }
              "
            ></b-form-input>
            <b-input-group-append>
              <b-button
                :disabled="!search"
                @click="
                  search = ''
                  currentPage = 1
                  searchPosts()
                "
                >Clear</b-button
              >
            </b-input-group-append>
          </b-input-group>
        </b-form-group>
      </b-col>
      <b-col md="6" class="my-1">
        <b-form-group label-cols-sm="3" label="Sort" class="mb-0">
          <b-input-group>
            <b-form-select
              v-model="sortBy"
              :options="sortOptions"
              @change="
                currentPage = 1
                searchPosts()
              "
            >
              <option slot="first" :value="null">-- none --</option>
            </b-form-select>
            <b-form-select
              slot="prepend"
              v-model="sortDesc"
              :disabled="!sortBy"
              @change="
                currentPage = 1
                searchPosts()
              "
            >
              <option :value="false">Asc</option>
              <option :value="true">Desc</option>
            </b-form-select>
          </b-input-group>
        </b-form-group>
      </b-col>
      <b-col md="6" class="my-1">
        <b-form-group label-cols-sm="3" label="Per page" class="mb-0">
          <b-form-select
            v-model="perPage"
            :options="pageOptions"
            @change="
              currentPage = 1
              searchPosts()
            "
          ></b-form-select>
        </b-form-group>
      </b-col>
    </b-row>
    <b-table
      :items="items"
      :fields="fields"
      :no-local-sorting="true"
      show-empty
      stacked="md"
      @sort-changed="sort"
    >
      <template v-slot:cell(title)="data">{{ data.value }}</template>
      <template v-slot:cell(author)="data">{{ data.value }}</template>
      <template v-slot:cell(date)="data">{{
        formatDate(Date.parse(data.value), 'yyyy-M-d')
      }}</template>
      <template v-slot:cell(views)="data">{{ data.value }}</template>
      <template v-slot:cell(read)="data">
        <nuxt-link :to="`/blog/${data.item.id}`" class="btn btn-primary btn-sm"
          >Read</nuxt-link
        >
      </template>
    </b-table>
    <b-row class="send-bottom">
      <b-col md="6" class="my-4">
        <b-pagination
          v-model="currentPage"
          :total-rows="totalRows"
          :per-page="perPage"
          class="my-0"
          @change="
            (newpage) => {
              currentPage = newpage
              searchPosts()
            }
          "
        ></b-pagination>
      </b-col>
    </b-row>
  </b-container>
</template>

<script lang="js">
import Vue from 'vue'
import { format } from 'date-fns'
import { validTypes } from '~/assets/config'
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
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
      items: [],
      fields: [
        {
          key: 'title',
          label: 'Title',
          sortable: true,
          sortDirection: 'desc'
        },
        {
          key: 'author',
          label: 'Author',
          sortable: true
        },
        {
          key: 'date',
          label: 'Date',
          sortable: true,
          class: 'text-center'
        },
        {
          key: 'views',
          label: 'Views',
          sortable: true,
          class: 'text-center'
        },
        {
          key: 'read',
          label: 'Read',
          sortable: false
        }
      ],
      totalRows: 0,
      currentPage: 1,
      perPage: 5,
      pageOptions: [5, 10, 15],
      sortBy: null,
      sortDesc: false,
      search: ''
    }
  },
  computed: {
    sortOptions() {
      // Create an options list from our fields
      return this.fields
        .filter(f => f.sortable)
        .map(f => {
          return { text: f.label, value: f.key }
        })
    }
  },
  mounted() {
    if (this.$route.query) {
      if (this.$route.query.phrase) this.search = this.$route.query.phrase
      if (this.$route.query.perpage)
        this.perPage = parseInt(this.$route.query.perpage)
      if (this.$route.query.currentpage)
        this.currentPage = parseInt(this.$route.query.currentpage)
      if (this.$route.query.sortdescending)
        this.sortDesc = this.$route.query.sortdescending === 'true'
      if (
        this.$route.query.sortby &&
        this.fields.some(field => field.key === this.$route.query.sortby)
      )
        this.sortBy = this.$route.query.sortby
    }
    this.searchPosts(this.currentPage)
  },
  methods: {
    sort(ctx) {
      this.sortBy = ctx.sortBy //   ==> Field key for sorting by (or null for no sorting)
      this.sortDesc = ctx.sortDesc // ==> true if sorting descending, false otherwise
      this.currentPage = 1
      this.searchPosts(this.currentPage)
    },
    updateCount() {
      this.$axios
        .get('/countPosts', {
          params: {
            searchterm: this.search,
            type: this.type,
            tags: [].join(',tags='),
            categories: [].join(',categories=')
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.count !== null) {
                this.totalRows = res.data.count
              } else {
                this.$toasted.global.error({
                  message: 'could not find count data'
                })
              }
            } else {
              this.$toasted.global.error({
                message: 'could not get data'
              })
            }
          } else {
            this.$toasted.global.error({
              message: `status code of ${res.status}`
            })
          }
        })
        .catch(err => {
          let message = `got error: ${err}`
          if (err.response && err.response.data) {
            message = err.response.data.message
          }
          this.$toasted.global.error({
            message
          })
        })
    },
    searchPosts() {
      this.updateCount()
      const sort = this.sortBy ? this.sortBy : this.sortOptions[0].value
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${encodeURIComponent(this.type)}",perpage:${
              this.perPage
            },page:${this.currentPage - 1},searchterm:"${encodeURIComponent(
              this.search
            )}",sort:"${encodeURIComponent(sort)}",ascending:${!this
              .sortDesc},tags:${JSON.stringify([])},categories:${JSON.stringify(
              []
            )},cache:${(!(
              this.$store.state.auth.user &&
              this.$store.state.auth.user.type === 'admin'
            )).toString()}){title views id author date}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                const posts = res.data.data.posts
                posts.forEach(post => {
                  Object.keys(post).forEach(key => {
                    if (typeof post[key] === 'string')
                      post[key] = decodeURIComponent(post[key])
                  })
                })
                this.items = posts
              } else if (res.data.errors) {
                this.$toasted.global.error({
                  message: `found errors: ${JSON.stringify(res.data.errors)}`
                })
              } else {
                this.$toasted.global.error({
                  message: 'could not find data or errors'
                })
              }
            } else {
              this.$toasted.global.error({
                message: 'could not get data'
              })
            }
          } else {
            this.$toasted.global.error({
              message: `status code of ${res.status}`
            })
          }
        })
        .catch(err => {
          let message = `got error: ${err}`
          if (err.response && err.response.data) {
            message = err.response.data.message
          }
          this.$toasted.global.error({
            message
          })
        })
    },
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    }
  },
  // @ts-ignore
  head() {
    const title = `Search ${this.type}`
    const description = `search for ${this.type}s, by name, views, etc`
    const image = `${seo.url}/icon.png`
    return {
      title,
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

<style lang="scss">
.main-container {
  min-height: 70vh;
}
.send-bottom {
  position: absolute;
  bottom: 0;
}
</style>
