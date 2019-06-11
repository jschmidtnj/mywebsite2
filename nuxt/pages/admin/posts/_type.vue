<template>
  <div class="container-fluid">
    <div id="admin-cards" class="container">
      <div class="row my-4">
        <div class="col-lg-6 my-2">
          <section class="card h-100 py-0">
            <div class="card-body">
              <b-form @submit="manageposts" @reset="resetposts">
                <span class="card-text">
                  <h2 class="mb-4">{{ mode }} Post</h2>
                  <b-form-group>
                    <label>Title</label>
                    <span>
                      <b-form-textarea
                        v-model="post.title"
                        type="text"
                        :state="!$v.post.title.$invalid"
                        class="form-control"
                        aria-describedby="titlefeedback"
                        placeholder="Enter title..."
                        rows="5"
                        max-rows="15"
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="titlefeedback"
                      :state="!$v.post.title.$invalid"
                    >
                      <div v-if="!$v.post.title.required">
                        title is required
                      </div>
                      <div v-else-if="!$v.post.title.minLength">
                        title must have at least
                        {{ $v.post.title.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Author</label>
                    <span>
                      <b-form-input
                        id="author"
                        v-model="post.author"
                        :state="!$v.post.author.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="authorfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="authorfeedback"
                      :state="!$v.post.author.$invalid"
                    >
                      <div v-if="!$v.post.author.required">
                        author is required
                      </div>
                      <div v-else-if="!$v.post.author.minLength">
                        author must have at least
                        {{ $v.post.author.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Content</label>
                    <span>
                      <b-form-input
                        id="content"
                        v-model="post.content"
                        :state="!$v.post.content.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="contentfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="contentfeedback"
                      :state="!$v.post.content.$invalid"
                    >
                      <div v-if="!$v.post.content.required">
                        content is required
                      </div>
                      <div v-else-if="!$v.post.content.minLength">
                        content must have at least
                        {{ $v.post.content.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Date</label>
                    <span>
                      <b-form-input
                        v-model="post.date"
                        type="date"
                        :state="!$v.post.date.$invalid"
                        class="form-control mb-2"
                        aria-describedby="datefeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="datefeedback"
                      :state="!$v.post.date.$invalid"
                    >
                      <div v-if="!$v.post.date.required">date is required</div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Case Number</label>
                    <span>
                      <b-form-input
                        v-model="post.casenumber"
                        type="number"
                        :state="!$v.post.casenumber.$invalid"
                        class="form-control mb-2"
                        aria-describedby="casenumberfeedback"
                        placeholder
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="casenumberfeedback"
                      :state="!$v.post.casenumber.$invalid"
                    >
                      <div v-if="!$v.post.casenumber.required">
                        case number is required
                      </div>
                      <div v-else-if="!$v.post.casenumber.integer">
                        case number must be an integer
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Hero Image</label>
                    <span>
                      <b-form-file
                        v-model="post.heroimage"
                        accept="image/*"
                        :state="!$v.post.heroimage.$invalid"
                        class="mb-2 form-control"
                        aria-describedby="heroimagefeedback"
                        placeholder="Choose an image..."
                        drop-placeholder="Drop image here..."
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="heroimagefeedback"
                      :state="!$v.post.heroimage.$invalid"
                    >
                      <div v-if="!$v.post.heroimage.required">
                        hero image is required
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <div
                    v-for="(imagevalue, index) in $v.post.images.$each.$iter"
                    id="fileselecter"
                    :key="`file-${index}`"
                  >
                    <b-form-group class="mb-2">
                      <label class="form-required">Image Name</label>
                      <span>
                        <b-form-input
                          v-model="post.images[index].name"
                          :state="!imagevalue.name.$invalid"
                          type="text"
                          class="form-control"
                          placeholder
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!imagevalue.name.$invalid"
                      >
                        <div v-if="!imagevalue.name.required">
                          image name is required
                        </div>
                        <div v-else-if="!imagevalue.name.minLength">
                          image name must have at least
                          {{ imagevalue.name.$params.minLength.min }} characters
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                    <b-form-group>
                      <label class="form-required">Image</label>
                      <span>
                        <b-form-file
                          v-model.trim="post.images[index].file"
                          accept="image/*"
                          :state="!imagevalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose an image..."
                          drop-placeholder="Drop image here..."
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!imagevalue.file.$invalid"
                      >
                        <div v-if="!imagevalue.file.required">
                          image is required
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                  </div>
                  <b-container>
                    <b-row>
                      <b-col>
                        <b-btn
                          variant="primary"
                          class="mr-2 mt-2"
                          @click="
                            post.images.push({
                              name: '',
                              file: null,
                              uploaded: false,
                              id: ''
                            })
                          "
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2"
                              style="max-width: 16px;"
                              icon="plus-circle"
                            /> </no-ssr
                          >Add
                        </b-btn>
                        <b-btn
                          variant="primary"
                          class="mr-2 mt-2"
                          :disabled="post.images.length === 1"
                          @click="post.images.pop()"
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2"
                              style="max-width: 16px;"
                              icon="times"
                            /> </no-ssr
                          >Remove
                        </b-btn>
                        <b-btn
                          variant="primary"
                          type="submit"
                          class="mr-2 mt-2"
                          :disabled="$v.post.$invalid"
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2"
                              style="max-width: 16px;"
                              icon="angle-double-right"
                            /> </no-ssr
                          >Sub
                        </b-btn>
                        <b-btn variant="primary" type="reset" class="mr-4 mt-2">
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2"
                              style="max-width: 16px;"
                              icon="times"
                            /> </no-ssr
                          >Clear
                        </b-btn>
                      </b-col>
                    </b-row>
                  </b-container>
                </span>
              </b-form>
            </div>
          </section>
        </div>
        <div class="col-lg-6 my-2">
          <section class="card h-100 py-0">
            <div class="card-body">
              <b-form @submit="searchposts" @reset="clearsearch">
                <span class="card-text">
                  <div
                    v-if="post.content !== ''"
                    id="content-rendered"
                    class="mb-4"
                  >
                    <h2 class="mb-4">Content</h2>
                    <vue-markdown
                      :source="post.content"
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
                <template slot="name" slot-scope="row">{{
                  row.value
                }}</template>
                <template slot="id" slot-scope="row">{{ row.value }}</template>
                <template slot="actions" slot-scope="row">
                  <b-button
                    size="sm"
                    class="mr-1"
                    @click="editPost(row.item)"
                    >{{ modetypes.edit }}</b-button
                  >
                  <b-button size="sm" @click="deletePost(row.item)">{{
                    modetypes.delete
                  }}</b-button>
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
import VueMarkdown from 'vue-markdown'
import Prism from 'prismjs'
/**
 * posts edit
 */
const modetypes = {
  add: 'Add',
  edit: 'Edit',
  delete: 'Delete'
}
export default Vue.extend({
  name: 'Posts',
  // @ts-ignore
  layout: 'admin',
  components: {
    VueMarkdown
  },
  mixins: [validationMixin],
  // @ts-ignore
  data() {
    return {
      validtypes: ['blog', 'project'],
      modetypes: modetypes,
      mode: modetypes.add,
      postid: null,
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
      post: {
        title: '',
        content: '',
        author: '',
        views: null,
        date: null,
        heroimage: null,
        images: []
      }
    }
  },
  // @ts-ignore
  validations: {
    search: {
      required,
      minLength: minLength(3)
    },
    post: {
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
      },
      heroimage: {
        file: {
          required
        }
      },
      images: {
        $each: {
          name: {
            required,
            minLength: minLength(3)
          },
          file: {
            required
          }
        }
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
        href: `${seo.url}/admin/posts`
      })
    }
    return {
      title: 'Admin Posts',
      links: links
    }
  },
  mounted() {
    if (!(this.type && this.validtypes.includes(this.type))) {
      this.$toasted.global.error({
        message: `invalid type given: ${this.type}`
      })
      // default to blog
      this.$router.push('/admin/blog')
    }
  },
  /* eslint-disable */
  methods: {
    updateCodeHighlighting() {
      this.$nextTick(() => {
        Prism.highlightAll()
      })
    },
    editPost(searchresult) {
      this.postid = searchresult._id
      console.log(`id: ${this.postid}`)

      // get images
      const getimages = thepost => {
        let getimagecount = 0
        let totalimages = 1 + thepost.images.length
        let cont = true
        const finishedGets = () => {
          this.$toasted.global.success({
            message: `edit ${this.type} with id ${this.postid}`
          })
        }
        this.$axios
          .get('/getPostPicture', {
            params: {
              type: this.type,
              postid: this.postid,
              hero: true
            }
          })
          .then(res => {
            if (!cont) return
            if (res.status == 200) {
              getimagecount++
              if (totalimages === getimagecount) {
                finishedGets()
              }
            } else {
              this.$toasted.global.error({
                message: `got status code of ${res.status} on image upload`
              })
            }
          })
          .catch(err => {
            this.$toasted.global.error({
              message: `got error on image get: ${err}`
            })
          })
        for (let i = 0; i < thepost.images.length; i++) {
          if (!cont) break
          this.$axios
            .get('/getPostPicture', {
              params: {
                type: this.type,
                postid: this.postid,
                imageid: thepost.images[i]
              }
            })
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                getimagecount++
                if (totalimages === getimagecount) {
                  finishedGets()
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on image upload`
                })
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on image get: ${err}`
              })
            })
        }
        this.post.heroimage = null
        this.post.images = []
        this.mode = this.modetypes.edit
        
      }

      // get post data first
      this.$axios
        .get('/graphql', {
          params: {
            query: `{post(type:"${this.type}",id:"${
              this.postid
            }"){title content id author views images}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const thepost: any = res.data.data.post
                getimages(thepost)
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
    deletePost(searchresult) {
      const id = searchresult._id
      console.log(`id: ${id}`)
      this.$axios
        .put('/graphql', {
          params: {
            query: `mutation{deletePost(type:"${this.type}",id:"${id}"){}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                this.searchresults.splice(
                  this.searchresults.indexOf(searchresult),
                  1
                )
                this.$toasted.global.success({
                  message: 'post deleted'
                })
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
    searchposts(evt) {
      evt.preventDefault()
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${this.type}",perpage:10,page:0,searchterm:"${
              this.search
            }",sort:"title",ascending:false){title id}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                this.searchresults = res.data.data.posts
                this.$toasted.global.success({
                  message: `found ${this.searchresults.length} result${
                    this.searchresults.length === 1 ? '' : 's'
                  }`
                })
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
    clearsearch(evt) {
      if (evt) evt.preventDefault()
      this.search = ''
      this.searchresults = []
    },
    resetposts(evt) {
      if (evt) evt.preventDefault()
      this.post = {
        title: '',
        content: '',
        author: '',
        views: null,
        date: null,
        heroimage: null,
        images: []
      }
      this.mode = this.modetypes.add
      this.postid = null
    },
    manageposts(evt) {
      evt.preventDefault()
      let postid = this.postid

      // upload image logic
      const uploadImages = () => {
        let cont = true
        let imageuploadcount = 0
        let totaluploads =
          (this.blog.heroimage.uploaded ? 1 : 0) +
          this.blog.images.filter(image => image.uploaded).length
        const uploadImage = (image, imageid, update) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', image)
          const endpoint = update ? '/updatePostPicture' : '/createPostPicture'
          this.$axios
            .put('/updatePostPicture', formData, {
              params: {
                type: this.type,
                postid: postid,
                imageid: imageid
              },
              headers: {
                'Content-Type': 'multipart/form-data'
              }
            })
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                imageuploadcount++
                if (totaluploads === imageuploadcount) {
                  this.$toasted.global.success({
                    message: `${this.mode}ed ${this.type} with id ${postid}`
                  })
                  this.mode = this.modetypes.add
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on image upload`
                })
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on image upload: ${err}`
              })
            })
        }
        if (!this.post.heroimage.uploaded) {
          this.blog.heroimage.file = new File(
            [this.blog.heroimage.file],
            'hero',
            {
              type: this.blog.heroimage.file.type
            }
          )
          uploadImage(
            this.post.heroimage.file,
            'hero',
            this.post.heroimage.update
          )
        }
        for (let i = 0; i < this.blog.images.length; i++) {
          if (!cont) break
          if (!this.blog.images[i].uploaded) {
            this.blog.images[i].file = new File(
              [this.blog.images[i].file],
              this.blog.images[i].name,
              {
                type: this.blog.images[i].file.type
              }
            )
            uploadImage(
              this.post.images[i].file,
              this.post.images[i].id,
              this.post.images[i].update
            )
          }
        }
      }

      // send to database logic (do this first)
      const postdata = Object.assign({}, this.post)
      postdata.heroimage = null
      postdata.images = []
      if (this.mode === this.modetypes.add) {
        this.$axios
          .post('/graphql', {
            params: {
              query: `mutation{addPost(type:"${this.type}",title:"${
                this.post.title
              }",content:"${this.post.content}",author:"${
                this.post.author
              }"){id}}`
            }
          })
          .then(res => {
            if (res.status === 200) {
              if (res.data) {
                if (res.data.data && res.data.data.post) {
                  postid = res.data.data.post.id
                  uploadImages()
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
      } else {
        this.$axios
          .put('/graphql', {
            params: {
              query: `mutation{updatePost(type:"${this.type}",id:"${
                this.postid
              }",title:"${this.post.title}",content:"${
                this.post.content
              }",author:"${this.post.author}"){}}`
            }
          })
          .then(res => {
            if (res.status === 200) {
              if (res.data) {
                if (res.data.data && res.data.data.post) {
                  uploadImages()
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
      }
    }
  }
})
</script>

<style lang="scss">
@import '~/node_modules/prismjs/themes/prism.css';
</style>
