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
                        type="text"
                        :state="!$v.post.content.$invalid"
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
                    class="sampleimage"
                    :src="post.heroimage.src"
                  ></b-img>
                  <b-form-group>
                    <label class="form-required">Hero Image</label>
                    <span>
                      <b-form-file
                        v-model="post.heroimage.file"
                        accept="image/jpeg, image/png"
                        :state="!$v.post.heroimage.$invalid"
                        class="mb-2 form-control"
                        aria-describedby="heroimagefeedback"
                        placeholder="Choose an image..."
                        drop-placeholder="Drop image here..."
                        @input="
                          post.heroimage.uploaded = false
                          updateImageSrc(post.heroimage)
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
                    class="sampleimage"
                    :src="post.tileimage.src"
                  ></b-img>
                  <b-form-group>
                    <label class="form-required">Tile Image</label>
                    <span>
                      <b-form-file
                        v-model="post.tileimage.file"
                        accept="image/jpeg, image/png"
                        :state="!$v.post.tileimage.$invalid"
                        class="mb-2 form-control"
                        aria-describedby="tileimagefeedback"
                        placeholder="Choose an image..."
                        drop-placeholder="Drop image here..."
                        @input="
                          post.tileimage.uploaded = false
                          updateImageSrc(post.tileimage)
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
                  <h4 class="mt-4">Images</h4>
                  <div
                    v-for="(imagevalue, index) in $v.post.images.$each.$iter"
                    :key="`image-${index}`"
                  >
                    <b-img
                      v-if="post.images[index].file && post.images[index].src"
                      class="sampleimage"
                      :src="post.images[index].src"
                    ></b-img>
                    <br />
                    <code
                      v-if="
                        post.images[index].file &&
                          post.images[index].name &&
                          post.images[index].width &&
                          post.images[index].height &&
                          post.images[index].id
                      "
                      >{{ getImageTag(post.images[index]) }}</code
                    >
                    <b-form-group class="mb-2">
                      <label class="form-required">Image Name</label>
                      <span>
                        <b-form-input
                          v-model="post.images[index].name"
                          :state="!imagevalue.name.$invalid"
                          type="text"
                          class="form-control"
                          placeholder="name"
                          @input="post.images[index].uploaded = false"
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
                          v-model="post.images[index].file"
                          accept="image/jpeg, image/png"
                          :state="!imagevalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose an image..."
                          drop-placeholder="Drop image here..."
                          @input="
                            post.images[index].uploaded = false
                            updateImageSrc(post.images[index])
                          "
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
                  <b-container class="mt-4">
                    <b-row>
                      <b-col>
                        <b-btn
                          variant="primary"
                          class="mr-2"
                          @click="
                            post.images.push({
                              name: '',
                              file: null,
                              uploaded: false,
                              id: createId(),
                              src: null,
                              width: null,
                              height: null
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
                          variant="primary"
                          class="mr-2"
                          :disabled="post.images.length === 0"
                          @click="removeImage"
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
                  <h4 class="mt-4">Gifs</h4>
                  <div
                    v-for="(gifvalue, index) in $v.post.gifs.$each.$iter"
                    :key="`gif-${index}`"
                  >
                    <b-img
                      v-if="post.gifs[index].file && post.gifs[index].src"
                      class="sampleimage"
                      :src="post.gifs[index].src"
                    ></b-img>
                    <br />
                    <code
                      v-if="
                        (post.gifs[index].file || post.gifs[index].uploaded) &&
                          post.gifs[index].name &&
                          post.gifs[index].width &&
                          post.gifs[index].height &&
                          post.gifs[index].id
                      "
                      >{{ getGifTag(post.gifs[index]) }}</code
                    >
                    <b-form-group class="mb-2">
                      <label class="form-required">Gif Name</label>
                      <span>
                        <b-form-input
                          v-model="post.gifs[index].name"
                          :state="!gifvalue.name.$invalid"
                          type="text"
                          class="form-control"
                          placeholder="name"
                          @input="post.gifs[index].uploaded = false"
                        />
                      </span>
                      <b-form-invalid-feedback :state="!gifvalue.name.$invalid">
                        <div v-if="!gifvalue.name.required">
                          gif name is required
                        </div>
                        <div v-else-if="!gifvalue.name.minLength">
                          gif name must have at least
                          {{ gifvalue.name.$params.minLength.min }} characters
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                    <b-form-group>
                      <label class="form-required">Gif</label>
                      <span>
                        <b-form-file
                          v-model="post.gifs[index].file"
                          accept="image/gif"
                          :state="!gifvalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose a gif..."
                          drop-placeholder="Drop gif here..."
                          @input="
                            post.gifs[index].uploaded = false
                            updateImageSrc(post.gifs[index])
                          "
                        />
                      </span>
                      <b-form-invalid-feedback :state="!gifvalue.file.$invalid">
                        <div v-if="!gifvalue.file.gotFile">
                          gif is required
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
                            post.gifs.push({
                              name: '',
                              file: null,
                              uploaded: false,
                              id: createId(),
                              src: null,
                              width: null,
                              height: null
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
                          variant="primary"
                          class="mr-2"
                          :disabled="post.gifs.length === 0"
                          @click="removeGif"
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
                  <h4 class="mt-4">Videos</h4>
                  <div
                    v-for="(videovalue, index) in $v.post.videos.$each.$iter"
                    :key="`video-${index}`"
                  >
                    <video
                      v-if="
                        post.videos[index].file &&
                          post.videos[index].src &&
                          post.videos[index].id
                      "
                      :ref="`video-source-${post.videos[index].id}`"
                      controls
                      autoplay
                      class="sampleimage"
                      :type="post.videos[index].type"
                      allowfullscreen
                      :src="post.videos[index].src"
                    ></video>
                    <br />
                    <code
                      v-if="
                        (post.videos[index].file ||
                          post.videos[index].uploaded) &&
                          post.videos[index].name &&
                          post.videos[index].width &&
                          post.videos[index].height &&
                          post.videos[index].id &&
                          post.videos[index].type
                      "
                      >{{ getVideoTag(post.videos[index]) }}</code
                    >
                    <b-form-group class="mb-2">
                      <label class="form-required">Video Name</label>
                      <span>
                        <b-form-input
                          v-model="post.videos[index].name"
                          :state="!videovalue.name.$invalid"
                          type="text"
                          class="form-control"
                          placeholder="name"
                          @input="post.videos[index].uploaded = false"
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!videovalue.name.$invalid"
                      >
                        <div v-if="!videovalue.name.required">
                          video name is required
                        </div>
                        <div v-else-if="!videovalue.name.minLength">
                          video name must have at least
                          {{ videovalue.name.$params.minLength.min }} characters
                        </div>
                      </b-form-invalid-feedback>
                    </b-form-group>
                    <b-form-group>
                      <label class="form-required">Video</label>
                      <span>
                        <b-form-file
                          v-model="post.videos[index].file"
                          accept="video/*"
                          :state="!videovalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose a video..."
                          drop-placeholder="Drop video here..."
                          @input="
                            post.videos[index].uploaded = false
                            updateVideoSrc(post.videos[index])
                          "
                        />
                      </span>
                      <b-form-invalid-feedback
                        :state="!videovalue.file.$invalid"
                      >
                        <div v-if="!videovalue.file.gotFile">
                          video is required
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
                            post.videos.push({
                              name: '',
                              file: null,
                              uploaded: false,
                              id: createId(),
                              src: null,
                              width: null,
                              height: null,
                              type: null
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
                          variant="primary"
                          class="mr-2"
                          :disabled="post.videos.length === 0"
                          @click="removeVideo"
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
                  <h4 class="mt-4">Files</h4>
                  <div
                    v-for="(filevalue, index) in $v.post.files.$each.$iter"
                    :key="`file-${index}`"
                  >
                    <code
                      v-if="
                        (post.files[index].file ||
                          post.files[index].uploaded) &&
                          post.images[index].name &&
                          post.images[index].id
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
                          accept="*"
                          :state="!filevalue.file.$invalid"
                          class="mb-2 form-control"
                          placeholder="Choose a file..."
                          drop-placeholder="Drop file here..."
                          @input="
                            post.files[post.files.length - 1].uploaded = false
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
                              id: createId()
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
                          variant="primary"
                          class="mr-2"
                          :disabled="post.files.length === 0"
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
                          variant="primary"
                          type="submit"
                          :disabled="$v.post.$invalid"
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
                        type="text"
                        :state="!$v.search.$invalid"
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
                    variant="primary"
                    type="submit"
                    class="mr-4"
                    :disabled="$v.search.$invalid"
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
                <template slot="date" slot-scope="row">{{
                  formatDate(row.value, 'M/D/YYYY')
                }}</template>
                <template slot="id" slot-scope="row">
                  <a :href="`/${type}/${row.value}`">{{ row.value }}</a>
                </template>
                <template slot="actions" slot-scope="row">
                  <b-button size="sm" class="mr-1" @click="editPost(row.item)"
                    >Edit</b-button
                  >
                  <b-button size="sm" @click="deletePost(row.item)"
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

<script lang="ts">
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
  staticstorageindexes
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
      modetypes: modetypes,
      mode: modetypes.add,
      postid: new ObjectID().toString(),
      search: '',
      searchresults: [],
      currentpage: 1,
      numperpage: 10,
      categoryOptions: options.categoryOptions,
      tagOptions: options.tagOptions,
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
        images: [],
        gifs: [],
        videos: [],
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
      },
      gifs: {
        $each: {
          name: {
            required,
            minLength: minLength(3)
          },
          file: {
            gotFile
          }
        }
      },
      videos: {
        $each: {
          name: {
            required,
            minLength: minLength(3)
          },
          file: {
            gotFile
          }
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
      title: title,
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
        this.type === 'blog' ? staticstorageindexes.blogimages : staticstorageindexes.projectimages
      }/${this.postid}/${image.id}/original" src="${
        cloudStorageURLs.posts
      }/${this.type === 'blog' ? staticstorageindexes.blogimages : staticstorageindexes.projectimages}/${this.postid}/${
        image.id}/blur" class="lazy img-fluid" alt="${
        image.name
      }" data-width="${image.width}" data-height="${image.height}">`
    },
    getGifTag(gif) {
      return `<img src="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.bloggifs : staticstorageindexes.projectgifs
      }/${this.postid}/${gif.id}" alt="${
        gif.name
      }" class="img-fluid" data-width="${gif.width}" data-height="${gif.height}">`
    },
    getVideoTag(video) {
      return `<video class="img-fluid" data-width="${video.width}" data-height="${
        video.height
      }" alt="${video.name}" controls allowfullscreen><source src="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? staticstorageindexes.blogvideos : staticstorageindexes.projectvideos
      }/${this.postid}/${video.id}" type="${video.type}" /></video>`
    },
    getFileTag(file) {
      return `<a href="${cloudStorageURLs.posts}/${
        this.type === 'blog' ? 'blogfiles' : 'projectfiles'
      }/${this.postid}/${file.id}"></a>`
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
      // @ts-ignore
      video.type = video.file.type
      const reader = new FileReader()
      reader.onload = e => {
        // @ts-ignore
        video.src = e.target.result
        this.$nextTick(() => {
          const videotag = this.$refs[`video-source-${video.id}`][0]
          console.log(videotag)
          videotag.load()
          videotag.oncanplay = () => {
            // @ts-ignore
            video.height = videotag.videoHeight
            // @ts-ignore
            video.width = videotag.videoWidth
            console.log(videotag)
            console.log(video.file)
            console.log('done2')
          }
        })
      }
      reader.readAsDataURL(video.file)
      console.log('done')
    },
    removeGif() {
      const removedGif = this.post.gifs[this.post.gifs.length - 1]
      const finished = () => {
        this.post.gifs.pop()
        this.$toasted.global.success({
          message: `removed gif ${removedGif.id}`
        })
      }
      if (this.mode === this.modetypes.add || !removedGif.uploaded) {
        finished()
      } else if (removedGif.name && removedGif.id && this.mode === this.modetypes.edit) {
        this.$axios
          .delete('/deletePostGifs', {
            data: {
              fileids: [
                removedGif.id
              ],
              postid: this.postid,
              type: this.type
            }
          })
          .then(res => {
            if (res.status == 200) {
              finished()
            } else {
              this.$toasted.global.error({
                message: `got status code of ${res.status} on gif delete`
              })
            }
          })
          .catch(err => {
            let message = `got error on gif delete: ${err}`
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
    removeVideo() {
      const removedVideo = this.post.videos[this.post.videos.length - 1]
      const finished = () => {
        this.post.videos.pop()
        this.$toasted.global.success({
          message: `removed video ${removedVideo.id}`
        })
      }
      if (this.mode === this.modetypes.add || !removedVideo.uploaded) {
        finished()
      } else if (removedVideo.name && removedVideo.id && this.mode === this.modetypes.edit) {
        this.$axios
          .delete('/deletePostVideos', {
            data: {
              fileids: [
                removedVideo.id
              ],
              postid: this.postid,
              type: this.type
            }
          })
          .then(res => {
            if (res.status == 200) {
              finished()
            } else {
              this.$toasted.global.error({
                message: `got status code of ${res.status} on video delete`
              })
            }
          })
          .catch(err => {
            let message = `got error on video delete: ${err}`
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
              type: this.type
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
    removeImage() {
      const removedImage = this.post.images[this.post.images.length - 1]
      const finished = () => {
        this.post.images.pop()
        this.$toasted.global.success({
          message: `removed image ${removedImage.id}`
        })
      }
      if (this.mode === this.modetypes.add || !removedImage.uploaded) {
        finished()
      } else if (removedImage.name && removedImage.id && this.mode === this.modetypes.edit) {
        this.$axios
          .delete('/deletePostPictures', {
            data: {
              imageids: [
                removedImage.id
              ],
              postid: this.postid,
              type: this.type
            }
          })
          .then(res => {
            if (res.status == 200) {
              finished()
            } else {
              this.$toasted.global.error({
                message: `got status code of ${res.status} on image delete`
              })
            }
          })
          .catch(err => {
            let message = `got error on image delete: ${err}`
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
        let getimagecount = 0
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
                this.type === 'blog' ? staticstorageindexes.blogimages : staticstorageindexes.projectimages
              }/${this.postid}/${thepost.heroimage.id}/original`,
              {
                responseType: 'blob'
              }
            )
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                if (res.data) {
                  thepost.heroimage = {
                    name: thepost.heroimage.name,
                    file: res.data,
                    uploaded: true,
                    id: thepost.heroimage.id,
                    src: null,
                    width: thepost.heroimage.width,
                    height: thepost.heroimage.height
                  }
                  this.updateImageSrc(thepost.heroimage)
                  gothero = true
                  if (
                    thepost.images.length === getimagecount &&
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
            thepost.images.length === getimagecount &&
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
                this.type === 'blog' ? staticstorageindexes.blogimages : staticstorageindexes.projectimages
              }/${this.postid}/${thepost.tileimage.id}/original`,
              {
                responseType: 'blob'
              }
            )
            .then(res => {
              if (!cont) return
              if (res.status == 200) {
                if (res.data) {
                  thepost.tileimage = {
                    name: thepost.tileimage.name,
                    file: res.data,
                    uploaded: true,
                    id: thepost.tileimage.id,
                    src: null,
                    width: thepost.tileimage.width,
                    height: thepost.tileimage.height
                  }
                  this.updateImageSrc(thepost.tileimage)
                  gottile = true
                  if (
                    thepost.images.length === getimagecount &&
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
            thepost.images.length === getimagecount &&
            gothero && !finished
          ) {
            finished = true
            finishedGets()
          }
        }
        if (thepost.images.length > 0) {
          for (let i = 0; i < thepost.images.length; i++) {
            if (!cont) break
            axios
              .get(
                `${cloudStorageURLs.posts}/${
                  this.type === 'blog' ? staticstorageindexes.blogimages : staticstorageindexes.projectimages
                }/${this.postid}/${thepost.images[i].id}/original`,
                {
                  responseType: 'blob'
                }
              )
              .then(res => {
                if (!cont) return
                if (res.status == 200) {
                  if (res.data) {
                    thepost.images[getimagecount] = {
                      id: thepost.images[getimagecount].id,
                      name: thepost.images[getimagecount].name,
                      uploaded: true,
                      file: res.data,
                      src: null,
                      width: thepost.images[getimagecount].width,
                      height: thepost.images[getimagecount].height
                    }
                    this.updateImageSrc(thepost.images[getimagecount])
                    getimagecount++
                    if (
                      thepost.images.length === getimagecount &&
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
        } else {
          if (
            gothero && gottile && !finished
          ) {
            finished = true
            finishedGets()
          }
        }
        if (thepost.gifs.length > 0) {
          for (let i = 0; i < thepost.gifs.length; i++) {
            thepost.gifs[i] = {
              id: thepost.gifs[i].id,
              name: thepost.gifs[i].name,
              uploaded: true,
              file: null,
              width: thepost.gifs[i].width,
              height: thepost.gifs[i].height
            }
          }
        }
        if (thepost.videos.length > 0) {
          for (let i = 0; i < thepost.videos.length; i++) {
            thepost.videos[i] = {
              id: thepost.videos[i].id,
              name: thepost.videos[i].name,
              uploaded: true,
              file: null,
              width: thepost.videos[i].width,
              height: thepost.videos[i].height,
              type: thepost.videos[i].type
            }
          }
        }
        if (thepost.files.length > 0) {
          for (let i = 0; i < thepost.files.length; i++) {
            thepost.files[i] = {
              id: thepost.files[i].id,
              name: thepost.files[i].name,
              uploaded: true,
              file: null
            }
          }
        }
        if (
          thepost.images.length === getimagecount &&
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
            )}",cache:false){title content id author views images{name id width height} gifs{name id width height} videos{name id width height type} heroimage{name id width height} tileimage{name id width height} caption comments files{name id} categories tags color}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.post) {
                const thepost: any = res.data.data.post
                console.log(res.data.data)
                Object.keys(thepost).forEach(key => {
                  if (typeof thepost[key] === 'string')
                    thepost[key] = decodeURIComponent(thepost[key]);
                })
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
        images: [],
        gifs: [],
        videos: [],
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

      // upload image logic
      const upload = () => {
        let cont = true
        let uploadcount = 0
        let imageuploads = this.post.images.filter(image => !image.uploaded)
        let gifuploads = this.post.gifs.filter(gif => !gif.uploaded)
        let videouploads = this.post.videos.filter(video => !video.uploaded)
        let fileuploads = this.post.files.filter(file => !file.uploaded)
        let totaluploads =
          (!this.post.heroimage.uploaded && this.post.heroimage.file ? 1 : 0) +
          (!this.post.tileimage.uploaded && this.post.tileimage.file ? 1 : 0) +
          imageuploads.length +
          gifuploads.length +
          videouploads.length +
          fileuploads.length
        let finished = false
        const successMessage = () => {
          this.$toasted.global.success({
            message: `${this.mode}ed ${this.type} with id ${postid}`
          })
          this.resetposts(evt)
        }
        const uploadImage = (image, imageid) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', image)
          this.$axios
            .put('/writePostPicture', formData, {
              params: {
                type: this.type,
                imageid: imageid,
                postid: this.postid
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
                  message: `got status code of ${res.status} on image upload`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on image upload: ${err}`
              })
              cont = false
            })
        }
        const uploadGif = (gif, gifid) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', gif)
          this.$axios
            .put('/writePostGif', formData, {
              params: {
                type: this.type,
                fileid: gifid,
                postid: this.postid
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
                  message: `got status code of ${res.status} on gif upload`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on gif upload: ${err}`
              })
              cont = false
            })
        }
        const uploadVideo = (video, videoid) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', video)
          this.$axios
            .put('/writePostVideo', formData, {
              params: {
                type: this.type,
                fileid: videoid,
                postid: this.postid
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
                  message: `got status code of ${res.status} on video upload`
                })
                cont = false
              }
            })
            .catch(err => {
              this.$toasted.global.error({
                message: `got error on video upload: ${err}`
              })
              cont = false
            })
        }
        const uploadFile = (file, fileid) => {
          if (!cont) return
          const formData = new FormData()
          formData.append('file', file)
          this.$axios
            .put('/writePostFile', formData, {
              params: {
                type: this.type,
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
              this.$toasted.global.error({
                message: `got error on file upload: ${err}`
              })
              cont = false
            })
        }
        let uploadinghero = false
        if (!this.post.heroimage.uploaded && this.post.heroimage.file) {
          uploadinghero = true
          this.post.heroimage.file = new File(
            [this.post.heroimage.file],
            'hero',
            {
              type: this.post.heroimage.file.type
            }
          )
          uploadImage(
            this.post.heroimage.file,
            this.post.heroimage.id
          )
        }
        let uploadingtile = false
        if (!this.post.tileimage.uploaded && this.post.tileimage.file) {
          uploadingtile = true
          this.post.tileimage.file = new File(
            [this.post.tileimage.file],
            'tile',
            {
              type: this.post.tileimage.file.type
            }
          )
          uploadImage(
            this.post.tileimage.file,
            this.post.tileimage.id
          )
        }
        if (imageuploads.length > 0) {
          for (let i = 0; i < imageuploads.length; i++) {
            imageuploads[i].file = new File(
              [imageuploads[i].file],
              imageuploads[i].name,
              {
                type: imageuploads[i].file.type
              }
            )
            uploadImage(
              imageuploads[i].file,
              imageuploads[i].id
            )
          }
        }
        if (gifuploads.length > 0) {
          for (let i = 0; i < gifuploads.length; i++) {
            gifuploads[i].file = new File(
              [gifuploads[i].file],
              gifuploads[i].name,
              {
                type: gifuploads[i].file.type
              }
            )
            uploadGif(
              gifuploads[i].file,
              gifuploads[i].id
            )
          }
        }
        if (videouploads.length > 0) {
          for (let i = 0; i < videouploads.length; i++) {
            videouploads[i].file = new File(
              [videouploads[i].file],
              videouploads[i].name,
              {
                type: videouploads[i].file.type
              }
            )
            uploadVideo(
              videouploads[i].file,
              videouploads[i].id
            )
          }
        }
        if (fileuploads.length > 0) {
          for (let i = 0; i < fileuploads.length; i++) {
            fileuploads[i].file = new File(
              [fileuploads[i].file],
              fileuploads[i].name,
              {
                type: fileuploads[i].file.type
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
          imageuploads.length === 0 &&
          gifuploads.length === 0 &&
          videouploads.length === 0 &&
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
              this.post.heroimage.file ? `id:"${encodeURIComponent(this.post.heroimage.id)}",name:"hero",height:${this.post.heroimage.height},width:${this.post.heroimage.width}` : ''
            }},tileimage:{${
              this.post.tileimage.file ? `id:"${encodeURIComponent(this.post.tileimage.id)}",name:"tile",height:${this.post.tileimage.height},width:${this.post.tileimage.width}` : ''
            }},images:[${
              this.post.images.map(image =>
                `{id:"${encodeURIComponent(image.id)}",name:"${encodeURIComponent(image.name)}",height:${image.height},width:${image.width}}`
              )
            }],gifs:[${
              this.post.gifs.map(gif =>
                `{id:"${encodeURIComponent(gif.id)}",name:"${encodeURIComponent(gif.name)}",height:${gif.height},width:${gif.width}}`
              )
            }],videos:[${
              this.post.videos.map(video =>
                `{id:"${encodeURIComponent(video.id)}",name:"${encodeURIComponent(video.name)}",height:${video.height},width:${video.width},type:"${video.type}"}`
              )
            }],files:[${
              this.post.files.map(file =>
                `{id:"${encodeURIComponent(file.id)}",name:"${encodeURIComponent(file.name)}"`
              )
            }],tags:${JSON.stringify(
              this.post.tags.map(tag => encodeURIComponent(tag))
            )},categories:${JSON.stringify(
              this.post.categories.map(category => encodeURIComponent(category))
            )}){id}}`
          })
          .then(res => {
            console.log(
              `images ${JSON.stringify(
                this.post.images.map(image =>
                  image.id
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
              this.post.heroimage.file ? `id:"${encodeURIComponent(this.post.heroimage.id)}",name:"hero",height:${this.post.heroimage.height},width:${this.post.heroimage.width}` : ''
            }},tileimage:{${
              this.post.tileimage.file ? `id:"${encodeURIComponent(this.post.tileimage.id)}",name:"tile",height:${this.post.tileimage.height},width:${this.post.tileimage.width}` : ''
            }},images:[${
              this.post.images.map(image =>
                `{id:"${encodeURIComponent(image.id)}",name:"${encodeURIComponent(image.name)}",height:${image.height},width:${image.width}}`
              )
            }],gifs:[${
              this.post.gifs.map(gif =>
                `{id:"${encodeURIComponent(gif.id)}",name:"${encodeURIComponent(gif.name)}",height:${gif.height},width:${gif.width}}`
              )
            }],videos:[${
              this.post.videos.map(video =>
                `{id:"${encodeURIComponent(video.id)}",name:"${encodeURIComponent(video.name)}",height:${video.height},width:${video.width},type:"${video.type}"}`
              )
            }],files:[${
              this.post.files.map(file =>
                `{id:"${encodeURIComponent(file.id)}",name:"${encodeURIComponent(file.name)}"`
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
