<template>
  <div class="container-fluid">
    <div id="admin-cards" class="container">
      <div class="row my-4">
        <div class="col-lg-6 my-2">
          <section class="card h-100 py-0">
            <div class="card-body">
              <b-form @submit="manageposts" @reset="resetposts">
                <span class="card-text">
                  <h2 class="mb-4">{{ mode }} {{ type }}</h2>
                  <b-form-group>
                    <label>Content</label>
                    <span>
                      <b-form-textarea
                        v-model="post.content"
                        :state="!$v.post.content.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="contentfeedback"
                        placeholder="Enter content..."
                        rows="5"
                        max-rows="15"
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
                    <label class="form-required">Author</label>
                    <span>
                      <b-form-input
                        id="author"
                        v-model="post.author"
                        :state="!$v.post.author.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="authorfeedback"
                        placeholder="author"
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
                    <label class="form-required">Title</label>
                    <span>
                      <b-form-input
                        id="title"
                        v-model="post.title"
                        :state="!$v.post.title.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="titlefeedback"
                        placeholder="title"
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
                    <label class="form-required">Caption</label>
                    <span>
                      <b-form-input
                        id="caption"
                        v-model="post.caption"
                        :state="!$v.post.caption.$invalid"
                        type="text"
                        class="form-control"
                        aria-describedby="captionfeedback"
                        placeholder="caption"
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="captionfeedback"
                      :state="!$v.post.caption.$invalid"
                    >
                      <div v-if="!$v.post.caption.required">
                        caption is required
                      </div>
                      <div v-else-if="!$v.post.caption.minLength">
                        caption must have at least
                        {{ $v.post.caption.$params.minLength.min }} characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Theme Color</label>
                    <span>
                      <no-ssr>
                        <color-picker
                          v-model="post.color"
                          aria-describedby="colorfeedback"
                        />
                      </no-ssr>
                    </span>
                    <b-form-invalid-feedback
                      id="colorfeedback"
                      :state="!$v.post.color.$invalid"
                    >
                      <div v-if="!$v.post.color.required">
                        color is required
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Categories</label>
                    <span>
                      <no-ssr>
                        <v-select
                          v-model="post.categories"
                          :options="categoryOptions"
                          :multiple="true"
                          :taggable="true"
                          aria-describedby="categoryfeedback"
                        ></v-select>
                      </no-ssr>
                    </span>
                    <b-form-invalid-feedback
                      id="categoryfeedback"
                      :state="!$v.post.categories.$invalid"
                    >
                      <div v-if="!$v.post.categories.required">
                        categories is required
                      </div>
                      <div v-else-if="!$v.post.categories.minLength">
                        categories must have at least
                        {{ $v.post.categories.$params.$each.minLength.min }}
                        characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-form-group>
                    <label class="form-required">Tags</label>
                    <span>
                      <no-ssr>
                        <v-select
                          v-model="post.tags"
                          :options="tagOptions"
                          :multiple="true"
                          :taggable="true"
                          aria-describedby="tagfeedback"
                        ></v-select>
                      </no-ssr>
                    </span>
                    <b-form-invalid-feedback
                      id="tagfeedback"
                      :state="!$v.post.tags.$invalid"
                    >
                      <div v-if="!$v.post.tags.required">tags is required</div>
                      <div v-else-if="!$v.post.tags.minLength">
                        tags must have at least
                        {{ $v.post.tags.$params.$each.minLength.min }}
                        characters
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <b-img
                    v-if="post.heroimage.file && post.heroimage.src"
                    :src="post.heroimage.src"
                    class="sampleimage"
                  ></b-img>
                  <b-form-group>
                    <label class="form-required">Hero Image</label>
                    <span>
                      <b-form-file
                        v-model="post.heroimage.file"
                        :accept="validimages.join(', ')"
                        :state="!$v.post.heroimage.$invalid"
                        class="mb-2 form-control"
                        aria-describedby="heroimagefeedback"
                        placeholder="Choose an image..."
                        drop-placeholder="Drop image here..."
                        @input="
                          post.heroimage.uploaded = false
                          updateFileSrc(post.heroimage)
                        "
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
                  <b-img
                    v-if="post.tileimage.file && post.tileimage.src"
                    :src="post.tileimage.src"
                    class="sampleimage"
                  ></b-img>
                  <b-form-group>
                    <label class="form-required">Tile Image</label>
                    <span>
                      <b-form-file
                        v-model="post.tileimage.file"
                        :accept="validimages.join(', ')"
                        :state="!$v.post.tileimage.$invalid"
                        class="mb-2 form-control"
                        aria-describedby="tileimagefeedback"
                        placeholder="Choose an image..."
                        drop-placeholder="Drop image here..."
                        @input="
                          post.tileimage.uploaded = false
                          updateFileSrc(post.tileimage)
                        "
                      />
                    </span>
                    <b-form-invalid-feedback
                      id="tileimagefeedback"
                      :state="!$v.post.tileimage.$invalid"
                    >
                      <div v-if="!$v.post.tileimage.required">
                        tile image is required
                      </div>
                    </b-form-invalid-feedback>
                  </b-form-group>
                  <h4 class="mt-4">Files</h4>
                  <div
                    v-for="(filevalue, index) in $v.post.files.$each.$iter"
                    :key="`file-${index}`"
                  >
                    <b-img
                      v-if="
                        post.files[index].src &&
                        post.files[index].type &&
                        checkImageType(post.files[index].type)
                      "
                      :src="post.files[index].src"
                      class="sampleimage"
                    ></b-img>
                    <video
                      v-else-if="
                        post.files[index].src &&
                        post.files[index].id &&
                        post.files[index].type &&
                        checkVideoType(post.files[index].type)
                      "
                      :ref="`video-source-${post.files[index].id}`"
                      :type="post.files[index].type"
                      :src="post.files[index].src"
                      controls
                      autoplay
                      class="sampleimage"
                      allowfullscreen
                    ></video>
                    <br />
                    <code
                      v-if="
                        post.files[index].type === 'image/gif' &&
                        (post.files[index].file ||
                          post.files[index].uploaded) &&
                        post.files[index].name &&
                        post.files[index].width &&
                        post.files[index].height &&
                        post.files[index].id
                      "
                      >{{ getGifTag(post.files[index]) }}</code
                    >
                    <code
                      v-else-if="
                        post.files[index].file &&
                        post.files[index].type &&
                        checkImageType(post.files[index].type) &&
                        post.files[index].name &&
                        post.files[index].width &&
                        post.files[index].height &&
                        post.files[index].id
                      "
                      >{{ getImageTag(post.files[index]) }}</code
                    >
                    <code
                      v-else-if="
                        (post.files[index].file ||
                          post.files[index].uploaded) &&
                        post.files[index].type &&
                        checkVideoType(post.files[index].type) &&
                        post.files[index].name &&
                        post.files[index].width &&
                        post.files[index].height &&
                        post.files[index].id &&
                        post.files[index].type
                      "
                      >{{ getVideoTag(post.files[index]) }}</code
                    >
                    <code
                      v-else-if="
                        (post.files[index].file ||
                          post.files[index].uploaded) &&
                        post.files[index].name &&
                        post.files[index].id
                      "
                      >{{ getFileTag(post.files[index]) }}</code
                    >
                    <b-form-group class="mb-2">
                      <label class="form-required">File Name</label>
                      <span>
                        <b-form-input
                          v-model="post.files[index].name"
                          :state="!filevalue.name.$invalid"
                          type="text"
                          class="form-control"
                          placeholder="name"
                          @input="post.files[index].uploaded = false"
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!filevalue.name.$invalid"
                      >
                        <div v-if="!filevalue.name.required">
                          file name is required
                        </div>
                        <div v-else-if="!filevalue.name.minLength">
                          file name must have at least
                          {{ filevalue.name.$params.minLength.min }} characters
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                    <b-form-group>
                      <label class="form-required">File</label>
                      <span>
                        <b-form-file
                          v-model="post.files[index].file"
                          :accept="validfiles.join(', ')"
                          :state="!filevalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose a file..."
                          drop-placeholder="Drop file here..."
                          @input="
                            post.files[index].uploaded = false
                            updateFileSrc(post.files[index])
                          "
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!filevalue.file.$invalid"
                      >
                        <div v-if="!filevalue.file.gotFile">
                          file is required
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                  </div>
                  <b-container class="mt-4">
                    <b-row>
                      <b-col>
                        <b-btn
                          variant="primary"
                          class="mr-2"
                          @click="
                            post.files.push({
                              name: '',
                              file: null,
                              uploaded: false,
                              id: createId(),
                              src: null,
                              width: null,
                              height: null,
                              type: null,
                            })
                          "
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2 arrow-size-edit"
                              icon="plus-circle"
                            /> </no-ssr
                          >Add
                        </b-btn>
                        <b-btn
                          :disabled="post.files.length === 0"
                          variant="primary"
                          class="mr-2"
                          @click="removeFile"
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2 arrow-size-edit"
                              icon="times"
                            /> </no-ssr
                          >Remove
                        </b-btn>
                      </b-col>
                    </b-row>
                  </b-container>
                  <b-container class="mt-4">
                    <b-row>
                      <b-col>
                        <b-btn
                          :disabled="$v.post.$invalid || submitting"
                          variant="primary"
                          type="submit"
                        >
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2 arrow-size-edit"
                              icon="angle-double-right"
                            /> </no-ssr
                          >Submit
                        </b-btn>
                        <b-btn variant="primary" type="reset" class="mr-4">
                          <no-ssr>
                            <font-awesome-icon
                              class="mr-2 arrow-size-edit"
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
                      @rendered="updateMarkdown"
                    />
                  </div>
                  <h2 class="mb-4">Search</h2>
                  <b-form-group>
                    <label class="form-required">Query</label>
                    <span>
                      <b-form-input
                        v-model="search"
                        :state="!$v.search.$invalid"
                        type="text"
                        class="form-control mb-2"
                        aria-describedby="searchfeedback"
                        placeholder="search..."
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
                    :disabled="$v.search.$invalid"
                    variant="primary"
                    type="submit"
                    class="mr-4"
                  >
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2 arrow-size-edit"
                        icon="angle-double-right"
                      /> </no-ssr
                    >Search
                  </b-btn>
                  <b-btn variant="primary" type="reset" class="mr-4">
                    <no-ssr>
                      <font-awesome-icon
                        class="mr-2 arrow-size-edit"
                        icon="times"
                      /> </no-ssr
                    >Clear
                  </b-btn>
                  <br />
                  <br />
                </span>
              </b-form>
              <b-table
                :items="searchresults"
                :fields="fields"
                :current-page="currentpage"
                :per-page="numperpage"
                show-empty
                stacked="md"
              >
                <template v-slot:cell(name)="data">{{ data.value }}</template>
                <template v-slot:cell(date)="data">{{
                  formatDate(data.value, 'M/D/YYYY')
                }}</template>
                <template v-slot:cell(id)="data">
                  <nuxt-link :to="`/${type}/${data.value}`">{{
                    data.value
                  }}</nuxt-link>
                </template>
                <template v-slot:cell(actions)="data">
                  <b-button size="sm" class="mr-1" @click="editPost(data.item)"
                    >Edit</b-button
                  >
                  <b-button size="sm" @click="deletePost(data.item)"
                    >Del</b-button
                  >
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

<script lang="js">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, minLength } from 'vuelidate/lib/validators'
import VueMarkdown from 'vue-markdown'
import Prism from 'prismjs'
import { format } from 'date-fns'
import uuid from 'uuid/v1'
import axios from 'axios'
import { Chrome } from 'vue-color'
import LazyLoad from 'vanilla-lazyload'
import { ObjectID } from 'bson'
import {
  cloudStorageURLs,
  validTypes,
  options,
  defaultColor,
  staticstorageindexes,
  validimages,
  validfiles,
  paths
} from '~/assets/config'
const gotFile = (_, vm) => vm.uploaded || vm.src !== null
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
const lazyLoadInstance = new LazyLoad({
  elements_selector: '.lazy'
})
/**
 * posts edit
 */
const modetypes = {
  add: 'Add',
  edit: 'Edit',
  delete: 'Delete'
}
const originalHero = {
  name: 'hero',
  uploaded: false,
  file: null,
  id: uuid(),
  src: null,
  width: null,
  height: null
}
const originalTile = Object.assign({}, originalHero)
originalTile.id = uuid()
originalTile.name = 'tile'
export default Vue.extend({
  name: 'Posts',
  // @ts-ignore
  layout: 'admin',
  components: {
    VueMarkdown,
    'color-picker': Chrome
  },
  mixins: [validationMixin],
  props: {
    type: {
      default: null,
      type: String,
      required: true,
      validator: val => validTypes.includes(String(val))
    }
  },
  // @ts-ignore
  data() {
    return {
      submitting: false,
      modetypes,
      mode: modetypes.add,
      postid: new ObjectID().toString(),
      search: '',
      searchresults: [],
      currentpage: 1,
      numperpage: 10,
      categoryOptions: options.categoryOptions,
      tagOptions: options.tagOptions,
      validimages,
      validfiles,
      paths,
      fields: [
        {
          key: 'title',
          label: 'Title',
          sortable: true
        },
        {
          key: 'date',
          label: 'Date',
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
        caption: '',
        color: defaultColor,
        author: '',
        tags: [],
        categories: [],
        heroimage: Object.assign({}, originalHero),
        tileimage: Object.assign({}, originalTile),
        files: []
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
        required,
        minLength: minLength(3)
      },
      caption: {
        required,
        minLength: minLength(3)
      },
      content: {
        required,
        minLength: minLength(10)
      },
      color: {
        required
      },
      heroimage: {
        file: {}
      },
      tileimage: {
        file: {
          required
        }
      },
      tags: {
        $each: {
          required
        }
      },
      categories: {
        $each: {
          required
        }
      },
      files: {
        $each: {
          name: {
            required,
            minLength: minLength(3)
          },
          file: {
            gotFile
          }
        }
      }
    }
  },
  // @ts-ignore
  head() {
    const title = `Admin Edit ${this.type}`
    const description = `admin page for editing ${this.type}s`
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
  },
  /* eslint-disable */
  methods: {
    checkImageType(type) {
      return /^image\/.*$/.test(type)
    },
    checkVideoType(type) {
      return /^video\/.*$/.test(type)
    },
    updateMarkdown() {
      this.$nextTick(() => {
        Prism.highlightAll()
        if (lazyLoadInstance) {
          console.log('update lazyload')
          lazyLoadInstance.update()
        }
      })
    },
    createId() {
      return uuid()
    },
    mongoidToDate(id) {
      return parseInt(id.substring(0, 8), 16) * 1000
    },
    formatDate(dateUTC, formatStr) {
      return format(dateUTC, formatStr)
    },
    getImageTag(image) {
      return `<img data-src="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
      }/${this.postid}/${image.id + this.paths.original}" src="${
        cloudStorageURLs.posts
      }/${this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles}/${this.postid}/${
        image.id + this.paths.blur}" class="lazy img-fluid" alt="${
        image.name
      }" data-width="${image.width}" data-height="${image.height}">`
    },
    getGifTag(gif) {
      return `<img data-src="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
      }/${this.postid}/${gif.id + this.paths.original}" placeholder-original="${
        cloudStorageURLs.posts
      }/${this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles}/${this.postid}/${
        gif.id + this.paths.placeholder + this.paths.original}" src="${
        cloudStorageURLs.posts
      }/${this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles}/${this.postid}/${
        gif.id + this.paths.placeholder + this.paths.blur}" class="lazy img-fluid gif" alt="${
        gif.name
      }" data-width="${gif.width}" data-height="${gif.height}">`
    },
    getVideoTag(video) {
      return `<video class="img-fluid" data-width="${video.width}" data-height="${
        video.height
      }" alt="${video.name}" controls allowfullscreen><source src="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
      }/${this.postid}/${video.id}#t=0.1" type="${video.type}" /></video>`
    },
    getFileTag(file) {
      return `<a href="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
      }/${this.postid}/${file.id}/original" target="_blank">download</a>`
    },
    updateFileSrc(file) {
      if (file.file) {
        file.type = file.file.type
        if (this.checkVideoType(file.type))
          this.updateVideoSrc(file)
        else
          this.updateImageSrc(file) 
      }
    },
    updateImageSrc(image) {
      console.log('start image src')
      if (!image.file) return
      const img = new Image()
      img.onload = () => {
        console.log('image loaded')
        image.width = img.width
        image.height = img.height
        console.log(`image width: ${image.width}, height: ${image.height}`)
      }
      const reader = new FileReader()
      reader.onload = e => {
        // @ts-ignore
        image.src = e.target.result
        img.src = image.src
      }
      reader.readAsDataURL(image.file)
      console.log('done')
    },
    updateVideoSrc(video) {
      if (!video.file) return
      const reader = new FileReader()
      reader.onload = e => {
        // @ts-ignore
        video.src = e.target.result
        this.$forceUpdate()
        this.$nextTick(() => {
          const videotag = this.$refs[`video-source-${video.id}`][0]
          videotag.load()
          videotag.oncanplay = () => {
            // @ts-ignore
            video.height = videotag.videoHeight
            // @ts-ignore
            video.width = videotag.videoWidth
          }
        })
      }
      reader.readAsDataURL(video.file)
      console.log('done')
    },
    removeFile() {
      const removedFile = this.post.files[this.post.files.length - 1]
      const finished = () => {
        this.post.files.pop()
        this.$toasted.global.success({
          message: `removed file ${removedFile.id}`
        })
      }
      if (this.mode === this.modetypes.add || !removedFile.uploaded) {
        finished()
      } else if (removedFile.name && removedFile.id && this.mode === this.modetypes.edit) {
        this.$axios
          .delete('/deletePostFiles', {
            data: {
              fileids: [
                removedFile.id
              ],
              postid: this.postid,
              posttype: this.type
            }
          })
          .then(res => {
            if (res.status == 200) {
              finished()
            } else {
              this.$toasted.global.error({
                message: `got status code of ${res.status} on file delete`
              })
            }
          })
          .catch(err => {
            let message = `got error on file delete: ${err}`
            if (err.response && err.response.data) {
              message = err.response.data.message
            }
            this.$toasted.global.error({
              message: message
            })
          })
      } else {
        this.$toasted.global.error({
          message: 'no name or id found, or mode type not edit'
        })
      }
    },
    editPost(searchresult) {
      this.postid = searchresult.id

      // get images
      const getimages = thepost => {
        let getfilecount = 0
        let gothero = false
        let gottile = false
        let cont = true
        let finished = false
        const finishedGets = () => {
          this.mode = this.modetypes.edit
          this.post = thepost
          this.$toasted.global.success({
            message: `edit ${this.type} with id ${this.postid}`
          })
        }
        if (thepost.heroimage !== null) {
          axios
            .get(
              `${cloudStorageURLs.posts}/${
                this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
              }/${this.postid}/${thepost.heroimage.id + this.paths.original}`,
              {
                responseType: 'blob'
              }
            )
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                if (res.data) {
                  thepost.heroimage.file = res.data
                  thepost.heroimage.uploaded = true
                  thepost.heroimage.src = null
                  this.updateFileSrc(thepost.heroimage)
                  gothero = true
                  if (
                    thepost.files.length === getfilecount &&
                    gottile && !finished
                  ) {
                    finished = true
                    finishedGets()
                  }
                } else {
                  this.$toasted.global.error({
                    message: 'could not get image data'
                  })
                  cont = false
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on image upload`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on hero image get: ${err}`
              })
              cont = false
            })
        } else {
          thepost.heroimage = Object.assign({}, originalHero)
          thepost.heroimage.id = this.createId()
          gothero = true
          if (
            thepost.files.length === getfilecount &&
            gottile && !finished
          ) {
            finished = true
            finishedGets()
          }
        }
        if (thepost.tileimage !== null) {
          axios
            .get(
              `${cloudStorageURLs.posts}/${
                this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
              }/${this.postid}/${thepost.tileimage.id + this.paths.original}`,
              {
                responseType: 'blob'
              }
            )
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                if (res.data) {
                  thepost.tileimage.uploaded = true
                  thepost.tileimage.file = res.data
                  thepost.tileimage.src = null
                  this.updateFileSrc(thepost.tileimage)
                  gottile = true
                  if (
                    thepost.files.length === getfilecount &&
                    gothero && !finished
                  ) {
                    finished = true
                    finishedGets()
                  }
                } else {
                  this.$toasted.global.error({
                    message: 'could not get image data'
                  })
                  cont = false
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on image download`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on tile image get: ${err}`
              })
              cont = false
            })
        } else {
          thepost.tileimage = Object.assign({}, originalHero)
          thepost.tileimage.id = this.createId()
          gottile = true
          if (
            thepost.files.length === getfilecount &&
            gothero && !finished
          ) {
            finished = true
            finishedGets()
          }
        }
        const getImageFile = filedata => {
          if (!cont) return
          axios
            .get(
              `${cloudStorageURLs.posts}/${
                this.type === 'blog' ? staticstorageindexes.blogfiles : staticstorageindexes.projectfiles
              }/${this.postid}/${filedata.id + this.paths.original}`,
              {
                responseType: 'blob'
              }
            )
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                if (res.data) {
                  filedata.file = res.data
                  thepost.src = null
                  this.updateFileSrc(filedata)
                  getfilecount++
                  if (
                    thepost.files.length === getfilecount &&
                    gothero && gottile && !finished
                  ) {
                    finished = true
                    finishedGets()
                  }
                } else {
                  this.$toasted.global.error({
                    message: 'could not get image data'
                  })
                  cont = false
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on image download`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on image get: ${err}`
              })
              cont = false
            })
        }
        if (thepost.files.length > 0) {
          for (let i = 0; i < thepost.files.length; i++) {
            if (!cont) break
            thepost.files[i].uploaded = true
            if (this.checkImageType(thepost.files[i].type) && thepost.files[i].type !== 'image/gif')
              getImageFile(thepost.files[i])
            else
              getfilecount++
          }
        } else {
          if (
            gothero && gottile && !finished
          ) {
            finished = true
            finishedGets()
          }
        }
        if (
          thepost.files.length === getfilecount &&
          gothero &&
          gottile &&
          !finished
        ) {
          finished = true
          finishedGets()
        }
      }
      // get post data first
      this.$axios
        .get('/graphql', {
          params: {
            query: `{post(type:"${encodeURIComponent(
              this.type
            )}",id:"${encodeURIComponent(
              this.postid
            )}",cache:false){title content id author views heroimage{name id width height type} tileimage{name id width height type} caption comments files{name id width height type} categories tags color}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const thepost = res.data.data.post
                console.log(res.data.data)
                Object.keys(thepost).forEach(key => {
                  if (typeof thepost[key] === 'string')
                    thepost[key] = decodeURIComponent(thepost[key]);
                })
                for (let i = 0; i < thepost.files.length; i++) {
                  Object.keys(thepost.files[i]).forEach(key => {
                    if (typeof thepost.files[i][key] === 'string')
                      thepost.files[i][key] = decodeURIComponent(thepost.files[i][key]);
                  })
                }
                thepost.tags = thepost.tags.map(tag => decodeURIComponent(tag))
                thepost.categories = thepost.categories.map(category => decodeURIComponent(category))
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
      const id = searchresult.id
      this.$axios
        .delete('/graphql', {
          params: {
            query: `mutation{deletePost(type:"${encodeURIComponent(
              this.type
            )}",id:"${encodeURIComponent(id)}"){id}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.deletePost) {
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
            query: `{posts(type:"${encodeURIComponent(
              this.type
            )}",perpage:10,page:0,searchterm:"${encodeURIComponent(
              this.search
            )}",sort:"title",ascending:false,tags:${JSON.stringify(
              []
            )},categories:${JSON.stringify([])},cache:false){title id}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                res.data.data.posts.map(
                  post => {
                    Object.keys(post).forEach(key => {
                      if (typeof post[key] === 'string')
                        post[key] = decodeURIComponent(post[key]);
                    })
                    post.date = this.mongoidToDate(post.id)
                  }
                )
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
        caption: '',
        color: defaultColor,
        author: '',
        heroimage: Object.assign({}, originalHero),
        tileimage: Object.assign({}, originalTile),
        files: [],
        tags: [],
        categories: []
      }
      this.post.heroimage.id = this.createId()
      this.post.tileimage.id = this.createId()
      this.mode = this.modetypes.add
      this.postid = new ObjectID().toString()
    },
    manageposts(evt) {
      evt.preventDefault()
      let postid = this.postid
      this.submitting = true

      // upload image logic
      const upload = () => {
        let cont = true
        let uploadcount = 0
        let fileuploads = this.post.files.filter(file => !file.uploaded)
        let totaluploads =
          (!this.post.heroimage.uploaded && this.post.heroimage.file ? 1 : 0) +
          (!this.post.tileimage.uploaded && this.post.tileimage.file ? 1 : 0) +
          fileuploads.length
        let finished = false
        const successMessage = () => {
          this.$toasted.global.success({
            message: `${this.mode}ed ${this.type} with id ${postid}`
          })
          this.submitting = false
          this.resetposts(evt)
        }
        const uploadFile = (file, fileid) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', file)
          this.$axios
            .put('/writePostFile', formData, {
              params: {
                posttype: this.type,
                filetype: file.type,
                postid: this.postid,
                fileid: fileid
              },
              headers: {
                'Content-Type': 'multipart/form-data'
              }
            })
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                uploadcount++
                if (totaluploads === uploadcount && !finished) {
                  finished = true
                  successMessage()
                }
              } else {
                this.$toasted.global.error({
                  message: `got status code of ${res.status} on file upload`
                })
                cont = false
              }
            })
            .catch(err => {
              let message = `got error: ${err}`
              if (err.response && err.response.data) {
                message = err.response.data.message
              }
              console.log(message)
              this.$toasted.global.error({
                message: message
              })
            })
        }
        let uploadinghero = false
        if (!this.post.heroimage.uploaded && this.post.heroimage.file && this.post.heroimage.type) {
          uploadinghero = true
          this.post.heroimage.file = new File(
            [this.post.heroimage.file],
            'hero',
            {
              type: this.post.heroimage.type
            }
          )
          uploadFile(
            this.post.heroimage.file,
            this.post.heroimage.id
          )
        }
        let uploadingtile = false
        if (!this.post.tileimage.uploaded && this.post.tileimage.file && this.post.tileimage.type) {
          uploadingtile = true
          this.post.tileimage.file = new File(
            [this.post.tileimage.file],
            'tile',
            {
              type: this.post.tileimage.type
            }
          )
          uploadFile(
            this.post.tileimage.file,
            this.post.tileimage.id
          )
        }
        if (fileuploads.length > 0) {
          for (let i = 0; i < fileuploads.length; i++) {
            fileuploads[i].file = new File(
              [fileuploads[i].file],
              fileuploads[i].name,
              {
                type: fileuploads[i].type
              }
            )
            uploadFile(
              fileuploads[i].file,
              fileuploads[i].id
            )
          }
        }
        if (
          !uploadinghero &&
          fileuploads.length === 0 &&
          !finished
        ) {
          finished = true
          successMessage()
        }
      }

      // send to database logic (do this first)
      const color = this.post.color.hex8
        ? this.post.color.hex8
        : this.post.color.toUpperCase()
      if (this.mode === this.modetypes.add) {
        this.$axios
          .post('/graphql', {
            query: `mutation{addPost(type:"${encodeURIComponent(
              this.type
            )}",id:"${encodeURIComponent(
              this.postid
            )}",title:"${encodeURIComponent(
              this.post.title
            )}",content:"${encodeURIComponent(
              this.post.content
            )}",color:"${encodeURIComponent(
              color
            )}",caption:"${encodeURIComponent(
              this.post.caption
            )}",author:"${encodeURIComponent(
              this.post.author
            )}",heroimage:{${
              this.post.heroimage.file ? `id:"${encodeURIComponent(this.post.heroimage.id)}",name:"hero",height:${this.post.heroimage.height},width:${this.post.heroimage.width},type:"${this.post.heroimage.type}"` : ''
            }},tileimage:{${
              this.post.tileimage.file ? `id:"${encodeURIComponent(this.post.tileimage.id)}",name:"tile",height:${this.post.tileimage.height},width:${this.post.tileimage.width},type:"${this.post.tileimage.type}"` : ''
            }},files:[${
              this.post.files.map(file =>
                `{id:"${encodeURIComponent(file.id)}",name:"${encodeURIComponent(file.name)}",height:${file.height ? file.height : 0},width:${file.width ? file.width : 0},type:"${file.type}"}`
              )
            }],tags:${JSON.stringify(
              this.post.tags.map(tag => encodeURIComponent(tag))
            )},categories:${JSON.stringify(
              this.post.categories.map(category => encodeURIComponent(category))
            )}){id}}`
          })
          .then(res => {
            console.log(
              `files ${JSON.stringify(
                this.post.files.map(file =>
                  file.id
                )
              )}`
            )
            if (res.status === 200) {
              if (res.data) {
                if (res.data.data && res.data.data.addPost) {
                  postid = res.data.data.addPost.id
                  upload()
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
            query: `mutation{updatePost(type:"${encodeURIComponent(
              this.type
            )}",id:"${encodeURIComponent(
              this.postid
            )}",title:"${encodeURIComponent(
              this.post.title
            )}",content:"${encodeURIComponent(
              this.post.content
            )}",color:"${encodeURIComponent(
              color
            )}",caption:"${encodeURIComponent(
              this.post.caption
            )}",author:"${encodeURIComponent(
              this.post.author
            )}",heroimage:{${
              this.post.heroimage.file ? `id:"${encodeURIComponent(this.post.heroimage.id)}",name:"hero",height:${this.post.heroimage.height},width:${this.post.heroimage.width},type:"${this.post.heroimage.type}"` : ''
            }},tileimage:{${
              this.post.tileimage.file ? `id:"${encodeURIComponent(this.post.tileimage.id)}",name:"tile",height:${this.post.tileimage.height},width:${this.post.tileimage.width},type:"${this.post.tileimage.type}"` : ''
            }},files:[${
              this.post.files.map(file =>
                `{id:"${encodeURIComponent(file.id)}",name:"${encodeURIComponent(file.name)}",height:${file.height ? file.height : 0},width:${file.width ? file.width : 0},type:"${file.type}"}`
              )
            }],tags:${JSON.stringify(
              this.post.tags.map(tag => encodeURIComponent(tag))
            )},categories:${JSON.stringify(
              this.post.categories.map(category => encodeURIComponent(category))
            )}){id}}`
          })
          .then(res => {
            if (res.status === 200) {
              if (res.data) {
                if (res.data.data && res.data.data.updatePost) {
                  upload()
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
.arrow-size-edit {
  font-size: 1rem;
}
.markdown {
  overflow: auto;
  max-height: 20rem;
}
.sampleimage {
  max-width: 200px;
}
</style>
