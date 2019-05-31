<template>
  <div class="container-fluid">
    <div id="admin-cards" class="container">
      <div class="row my-4">
        <div class="col-lg-6 my-2">
          <section class="card h-100 py-0">
            <div class="card-body">
              <b-form @submit="manageblogs" @reset="resetblogs">
                <span class="card-text">
                  <h2 class="mb-4">{{ mode }} Blog</h2>
                  <b-form-group>
                    <label>Title</label>
                    <span>
                      <b-form-textarea
                        v-model="blog.title"
                        type="text"
                        :state="!$v.blog.title.$invalid"
                        class="form-control"
                        aria-describedby="titlefeedback"
                        placeholder="Enter title..."
                        rows="5"
                        max-rows="15"
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="titlefeedback"
                      :state="!$v.blog.title.$invalid"
                    >
                      <div v-if="!$v.blog.title.required">
                        title is required
                      </div>
                      <div v-else-if="!$v.blog.title.minLength">
                        title must have at least
                        {{ $v.blog.title.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Author</label>
                    <span>
                      <b-form-input
                        id="author"
                        v-model="blog.author"
                        :state="!$v.blog.author.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="authorfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="authorfeedback"
                      :state="!$v.blog.author.$invalid"
                    >
                      <div v-if="!$v.blog.author.required">
                        author is required
                      </div>
                      <div v-else-if="!$v.blog.author.minLength">
                        author must have at least
                        {{ $v.blog.author.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Content</label>
                    <span>
                      <b-form-input
                        id="content"
                        v-model="blog.content"
                        :state="!$v.blog.content.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="contentfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="contentfeedback"
                      :state="!$v.blog.content.$invalid"
                    >
                      <div v-if="!$v.blog.content.required">
                        content is required
                      </div>
                      <div v-else-if="!$v.blog.content.minLength">
                        content must have at least
                        {{ $v.blog.content.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Date</label>
                    <span>
                      <b-form-input
                        v-model="blog.date"
                        type="date"
                        :state="!$v.blog.date.$invalid"
                        class="form-control mb-2"
                        aria-describedby="datefeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="datefeedback"
                      :state="!$v.blog.date.$invalid"
                    >
                      <div v-if="!$v.blog.date.required">date is required</div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Case Number</label>
                    <span>
                      <b-form-input
                        v-model="blog.casenumber"
                        type="number"
                        :state="!$v.blog.casenumber.$invalid"
                        class="form-control mb-2"
                        aria-describedby="casenumberfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="casenumberfeedback"
                      :state="!$v.blog.casenumber.$invalid"
                    >
                      <div v-if="!$v.blog.casenumber.required">
                        case number is required
                      </div>
                      <div v-else-if="!$v.blog.casenumber.integer">
                        case number must be an integer
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-btn
                    variant="primary"
                    type="submit"
                    class="mr-4"
                    :disabled="$v.blog.$invalid"
                  >
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2"
                        style="max-width: 16px;"
                        icon="angle-double-right"
                      /> </no-ssr
                    >Submit
                  </b-btn>
                  <b-btn variant="primary" type="reset" class="mr-4">
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2"
                        style="max-width: 16px;"
                        icon="times"
                      /> </no-ssr
                    >Clear
                  </b-btn>
                </span>
              </b-form>
            </div>
          </section>
        </div>
        <div class="col-lg-6 my-2">
          <section class="card h-100 py-0">
            <div class="card-body">
              <b-form @submit="searchblogs" @reset="clearsearch">
                <span class="card-text">
                  <div
                    v-if="blog.content !== ''"
                    id="content-rendered"
                    class="mb-4"
                  >
                    <h2 class="mb-4">Content</h2>
                    <vue-markdown
                      :source="blog.content"
                      class="mb-4 markdown"
                      @rendered="updateCodeHighlighting"
                    />
                  </div>
                  <h2 class="mb-4">Search</h2>
                  <b-form-group>
                    <label class="form-required">Query</label>
                    <span>
                      <b-form-input
                        v-model="search"
                        type="text"
                        :state="!$v.search.$invalid"
                        class="form-control mb-2"
                        aria-describedby="searchfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="searchfeedback"
                      :state="!$v.search.$invalid"
                    >
                      <div v-if="!$v.search.required">query is required</div>
                      <div v-else-if="!$v.search.minLength">
                        query must have at least
                        {{ $v.search.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-btn
                    variant="primary"
                    type="submit"
                    class="mr-4"
                    :disabled="$v.search.$invalid"
                  >
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2"
                        style="max-width: 16px;"
                        icon="angle-double-right"
                      /> </no-ssr
                    >Search
                  </b-btn>
                  <b-btn variant="primary" type="reset" class="mr-4">
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2"
                        style="max-width: 16px;"
                        icon="times"
                      /> </no-ssr
                    >Clear
                  </b-btn>
                  <br />
                  <br />
                </span>
              </b-form>
              <b-table
                show-empty
                stacked="md"
                :items="searchresults"
                :fields="fields"
                :current-page="currentpage"
                :per-page="numperpage"
              >
                <template slot="name" slot-scope="row">
                  {{ row.value }}
                </template>
                <template slot="id" slot-scope="row">
                  {{ row.value }}
                </template>
                <template slot="actions" slot-scope="row">
                  <b-button size="sm" class="mr-1" @click="editBlog(row.item)">
                    {{ modetypes.edit }}
                  </b-button>
                  <b-button size="sm" @click="deleteBlog(row.item)">
                    {{ modetypes.delete }}
                  </b-button>
                </template>
              </b-table>
              <b-row class="mb-2">
                <b-col md="6" class="my-1">
                  <b-pagination
                    v-model="currentpage"
                    :total-rows="searchresults.length"
                    :per-page="numperpage"
                    class="my-0"
                  />
                </b-col>
              </b-row>
            </div>
          </section>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, minLength, integer } from 'vuelidate/lib/validators'
import { codes } from '~/assets/config'
import VueMarkdown from 'vue-markdown'
import Prism from 'prismjs'
/**
 * blogs edit
 */
const modetypes = {
  add: 'Add',
  edit: 'Edit',
  delete: 'Delete'
}
export default Vue.extend({
  name: 'Blogs',
  // @ts-ignore
  layout: 'admin',
  components: {
    VueMarkdown
  },
  mixins: [validationMixin],
  // @ts-ignore
  data() {
    return {
      modetypes: modetypes,
      mode: modetypes.add,
      blogid: null,
      search: '',
      searchresults: [],
      currentpage: 1,
      numperpage: 10,
      fields: [
        {
          key: 'title',
          label: 'Title',
          sortable: true
        },
        {
          key: 'id',
          label: 'ID',
          sortable: true
        },
        {
          key: 'actions',
          label: 'Actions',
          sortable: false
        }
      ],
      blog: {
        title: '',
        content: '',
        author: '',
        views: null,
        date: null
      }
    }
  },
  // @ts-ignore
  validations: {
    search: {
      required,
      minLength: minLength(3)
    },
    blog: {
      title: {
        required,
        minLength: minLength(3)
      },
      author: {
        minLength: minLength(3)
      },
      content: {
        required,
        minLength: minLength(10)
      },
      views: {
        required,
        integer
      }
    }
  },
  head() {
    // @ts-ignore
    const seo = JSON.parse(process.env.seoconfig)
    const links: any = []
    if (seo) {
      links.push({
        rel: 'canonical',
        href: `${seo.url}/admin/blogs`
      })
    }
    return {
      title: 'Admin Blogs',
      links: links
    }
  },
  mounted() {},
  methods: {
    updateCodeHighlighting() {
      this.$nextTick(() => {
        Prism.highlightAll()
      })
    },
    editBlog(searchresult) {
      this.$axios
        .get()
        .then(res => {
          this.blog = res.data
          this.mode = this.modetypes.edit
          this.$toasted.global.success({
            message: 'edit blog'
          })
        })
        .catch(err => {
          this.$toasted.global.error({
            message: err
          })
        })
    },
    deleteBlog(searchresult) {
      this.$axios
        .put()
        .then(res => {
          this.searchresults.splice(this.searchresults.indexOf(searchresult), 1)
          this.$toasted.global.success({
            message: 'blog deleted'
          })
        })
        .catch(err => {
          this.$toasted.global.error({
            message: err
          })
        })
    },
    searchblogs(evt) {
      evt.preventDefault()
      const searchfunction = null
      searchfunction({
        request: {
          _source: ['firstname', 'lastname'],
          query: {
            multi_match: {
              query: this.search,
              type: 'best_fields',
              tie_breaker: 0.3
            }
          }
        }
      })
        .then(res => {
          let message = ''
          let code = 0
          switch (res.data.code) {
            case codes.success:
              this.searchresults = res.data.message.hits.hits
              message = `found ${this.searchresults.length} result${
                this.searchresults.length === 1 ? '' : 's'
              }`
              code =
                this.searchresults.length !== 0 ? codes.success : codes.error
              break
            case codes.error:
              message = `got error ${res.data.message}`
              code = codes.error
              break
            case codes.unauthorized:
              message = `unauthorized: you must be signed in to search`
              code = codes.unauthorized
              break
            default:
              message = `warning: ${res.data.message}`
              code = codes.warning
              break
          }
          this.$store.commit('notifications/addNotification', {
            code: code,
            message: message
          })
        })
        .catch(err => {
          this.$store.commit('notifications/addNotification', {
            code: codes.error,
            message: JSON.stringify(err)
          })
        })
    },
    clearsearch(evt) {
      if (evt) evt.preventDefault()
      this.search = ''
      this.searchresults = []
    },
    resetblogs(evt) {
      if (evt) evt.preventDefault()
      this.blog = {
        title: '',
        content: '',
        author: '',
        views: null,
        date: null
      }
      this.mode = this.modetypes.add
      this.blogid = null
    },
    manageblogs(evt) {
      evt.preventDefault()
      const blogdata = Object.assign({}, this.blog)
      blogdata.updated = new Date().getTime()
      blogdata.articles = blogdata.articles.map(result => result.id)
      blogdata.yearjoined = parseInt(blogdata.yearjoined)
      blogdata.baradmissions.year = parseInt(blogdata.baradmissions.year)
      delete blogdata.image
      const successmessage = () => {
        this.$store.commit('notifications/addNotification', {
          code: codes.success,
          message: `added blog ${blogdata.firstname} ${blogdata.lastname}`
        })
      }
      if (this.mode === this.modetypes.add) {
        blogdata.created = new Date().getTime()
        const Firestore = null
        Firestore.collection('blogs')
          .add(blogdata)
          .then(docRef => {
            this.blog.image = new File([this.blog.image], docRef.id, {
              type: this.blog.image.type
            })
            const Storage = null
            Storage.ref('blogimages')
              .child(docRef.id)
              .put(this.blog.image)
              .then(storagesnap => {
                successmessage()
                // this.resetblogs()
              })
              .catch(err => {
                this.$store.commit('notifications/addNotification', {
                  code: codes.error,
                  message: JSON.stringify(err)
                })
              })
          })
          .catch(err => {
            this.$store.commit('notifications/addNotification', {
              code: codes.error,
              message: JSON.stringify(err)
            })
          })
      } else {
        const Firestore = null
        Firestore.collection('blogimages')
          .doc(this.blogid)
          .set(blogdata)
          .then(() => {
            if (this.blog.image) {
              this.blog.image = new File([this.blog.image], this.blogid, {
                type: this.blog.image.type
              })
              const Storage = null
              Storage.ref('blogimages')
                .child(this.blogid)
                .put(this.blog.image)
                .then(storagesnap => {
                  successmessage()
                  this.resetblogs()
                })
                .catch(err => {
                  this.$store.commit('notifications/addNotification', {
                    code: codes.error,
                    message: JSON.stringify(err)
                  })
                })
            } else {
              successmessage()
              this.resetblogs()
            }
          })
          .catch(err => {
            this.$store.commit('notifications/addNotification', {
              code: codes.error,
              message: JSON.stringify(err)
            })
          })
          .then(() => {
            this.blogid = null
            this.mode = this.modetypes.add
          })
      }
    }
  }
})
</script>

<style lang="scss"></style>
