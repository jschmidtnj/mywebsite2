<template>
  <b-card title="Create" footer-tag="footer">
    <b-form @submit="createLink">
      <b-form-group
        id="link-group"
        label="Link:"
        label-for="link"
        description="Insert url here"
      >
        <b-form-input
          id="link"
          v-model="form.link"
          type="text"
          :state="!$v.form.link.$invalid"
          placeholder="Enter url"
          aria-describedby="linkfeedback"
        ></b-form-input>
        <b-form-invalid-feedback
          id="linkfeedback"
          :state="!$v.form.link.$invalid"
        >
          <div v-if="!$v.form.link.required">url is required</div>
          <div v-else-if="!$v.form.link.url">url is invalid</div>
        </b-form-invalid-feedback>
      </b-form-group>
      <br />
      <b-button
        variant="primary"
        type="submit"
        class="mt-4"
        :disabled="$v.form.$invalid"
        >Submit</b-button
      >
    </b-form>
    <p slot="footer">
      By clicking submit you aggree to the
      <a :href="`${mainurl}/privacy`">privacy policy</a>. This site is protected
      by reCAPTCHA and the Google
      <a href="https://policies.google.com/privacy">Privacy Policy</a> and
      <a href="https://policies.google.com/terms">Terms of Service</a> apply.
    </p>
    <p v-if="shortlink">{{ `${currenturl}/${shortlink}` }}</p>
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, url } from 'vuelidate/lib/validators'
const mainurl = process.env.mainurl ? process.env.mainurl : ''
// @ts-ignore
const currenturl = JSON.parse(process.env.seoconfig).url
export default Vue.extend({
  name: 'Home',
  mixins: [validationMixin],
  data() {
    return {
      form: {
        link: ''
      },
      shortlink: null,
      mainurl: mainurl,
      currenturl: currenturl
    }
  },
  // @ts-ignore
  validations: {
    form: {
      link: {
        required,
        url
      }
    }
  },
  // @ts-ignore
  head() {
    return {
      title: 'Home'
    }
  },
  mounted() {
    /* eslint-disable */
    console.log('hello world!')
  },
  methods: {
    createLink(evt) {
      evt.preventDefault()
      this.$recaptcha('login')
        .then(recaptchatoken => {
          console.log(`got recaptcha token ${recaptchatoken}`)
          if (this.$store.state.auth && this.$store.state.auth.loggedIn) {
            console.log('user logged in')
            this.$axios
              .post(
                '/graphql',
                {},
                {
                  params: {
                    query: `mutation{addShortlink(link:"${encodeURIComponent(
                      this.form.link
                    )}",recaptcha:"${encodeURIComponent(recaptchatoken)}"){id}}`
                  }
                }
              )
              .then(res => {
                if (res.status === 200) {
                  if (res.data) {
                    if (res.data.data && res.data.data.addShortlink) {
                      this.shortlink = res.data.data.addShortlink.id
                      this.$toasted.global.success({
                        message: 'created short link'
                      })
                      const newShortlinks = [...this.$store.state.auth.user.shortlinks]
                      newShortlinks.push(this.shortlink)
                      this.$store.commit('auth/updateUser', {
                        path: 'shortlinks',
                        value: newShortlinks
                      })
                    } else if (res.data.errors) {
                      this.$toasted.global.error({
                        message: `found errors: ${JSON.stringify(
                          res.data.errors
                        )}`
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
                  message: message
                })
              })
          } else {
            this.$axios
              .post('/createShortLink', {
                link: encodeURIComponent(this.form.link),
                recaptcha: recaptchatoken
              })
              .then(res => {
                if (res.status === 200) {
                  if (res.data) {
                    this.shortlink = res.data.id
                    console.log(`short link ${this.shortlink}`)
                    this.$toasted.global.success({
                      message: 'created short link'
                    })
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
                  message: message
                })
              })
          }
        })
        .catch(err => {
          this.$toasted.global.error({
            message: `got error with recaptcha ${err}`
          })
        })
    }
  }
})
</script>

<style lang="scss" scoped>
</style>
