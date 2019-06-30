<template>
  <b-container fluid>
    <b-row>
      <b-col md="6" class="my-1">
        <b-form-group label-cols-sm="3" label="search" class="mb-0">
          <b-input-group>
            <b-form-input
              v-model="search"
              placeholder="Type to Search"
              @keyup.enter.native="
                evt => {
                  evt.preventDefault()
                  pageNum = 1
                  searchPosts(currentPage)
                }
              "
            ></b-form-input>
            <b-input-group-append>
              <b-button
                :disabled="!search"
                @click="
                  search = ''
                  pageNum = 1
                  searchPosts(currentPage)
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
                searchPosts(currentPage)
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
                searchPosts(currentPage)
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
              searchPosts(currentPage)
            "
          ></b-form-select>
        </b-form-group>
      </b-col>
    </b-row>
    <b-table show-empty stacked="md" :items="items" :fields="fields">
      <template slot="title" slot-scope="row">{{ row.value }}</template>
      <template slot="author" slot-scope="row">{{ row.value }}</template>
      <template slot="date" slot-scope="row">{{
        formatDate(row.value, 'M/D/YYYY')
      }}</template>
      <template slot="views" slot-scope="row">{{ row.value }}</template>
      <template slot="read" slot-scope="row">
        <a class="btn btn-primary btn-sm" :href="`/blog?id=${row.item.id}`"
          >Read</a
        >
      </template>
    </b-table>
    <b-row>
      <b-col md="6" class="my-1">
        <b-pagination
          v-model="currentPage"
          :total-rows="totalRows"
          :per-page="perPage"
          class="my-0"
          @change="searchPosts"
        ></b-pagination>
      </b-col>
    </b-row>
  </b-container>
</template>

<script lang="ts">
import Vue from 'vue'
import { format } from 'date-fns'
import { validTypes } from '~/assets/config'
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
      sortDirection: 'asc',
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
    this.searchPosts(this.currentPage)
  },
  methods: {
    updateCount() {
      this.$axios
        .get('/countPosts', {
          params: {
            searchterm: this.search,
            type: this.type
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
          this.$toasted.global.error({
            message: `got error: ${err}`
          })
        })
    },
    searchPosts(pageNum) {
      this.updateCount()
      const sort = this.sortBy ? this.sortBy : this.sortOptions[0].value
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${this.type}",perpage:${
              this.perPage
            },page:${pageNum - 1},searchterm:"${
              this.search
            }",sort:"${sort}",ascending:${!this
              .sortDesc}){title views id author date}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                this.items = res.data.data.posts
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
          this.$toasted.global.error({
            message: err
          })
        })
    },
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    }
  }
})
</script>

<style lang="scss"></style>
