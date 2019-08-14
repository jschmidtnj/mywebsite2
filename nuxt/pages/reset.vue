<template>
  <div>
    <loading v-if="loading" />
    <b-card v-else title="Reset" footer-tag="footer">
      <b-form v-if="sendEmailMode" @submit="sendResetEmail">
        <b-form-group
          id="email-address-group"
          label="Email address:"
          label-for="email-address"
          description="Please enter your account email address"
        >
          <b-form-input
            id="email-address"
            v-model="emailform.email"
            type="text"
            autocomplete="off"
            :state="!$v.emailform.email.$invalid"
            placeholder="Enter email"
            aria-describedby="emailfeedback"
          ></b-form-input>
          <b-form-invalid-feedback
            id="emailfeedback"
            :state="!$v.emailform.email.$invalid"
          >
            <div v-if="!$v.emailform.email.required">email is required</div>
            <div v-else-if="!$v.emailform.email.email">
              email is invalid
            </div>
          </b-form-invalid-feedback>
        </b-form-group>
        <b-button
          variant="primary"
          type="submit"
          class="mt-4"
          :disabled="$v.emailform.$invalid"
          >Submit</b-button
        >
      </b-form>
      <b-form v-else @submit="resetPassword">
        <b-form-group
          id="password-group"
          label="Password:"
          label-for="password"
          description="Password must be at least 8 characters long, with a number, capital letter and special character"
        >
          <b-form-input
            id="password"
            v-model="passwordform.password"
            type="password"
            :state="!$v.passwordform.password.$invalid"
            placeholder="Enter password"
            aria-describedby="passwordfeedback"
          ></b-form-input>
          <b-form-invalid-feedback
            id="passwordfeedback"
            :state="!$v.passwordform.password.$invalid"
          >
            <div v-if="!$v.passwordform.password.required">
              password is required
            </div>
            <div v-else-if="!$v.passwordform.password.validPassword">
              password is invalid
            </div>
          </b-form-invalid-feedback>
        </b-form-group>
        <b-button
          variant="primary"
          type="submit"
          class="mt-4"
          :disabled="$v.passwordform.$invalid"
          >Submit</b-button
        >
      </b-form>
      <p slot="footer">
        By clicking submit you aggree to the
        <a href="/privacy">privacy policy</a>. This site is protected by
        reCAPTCHA and the Google
        <a href="https://policies.google.com/privacy">Privacy Policy</a> and
        <a href="https://policies.google.com/terms">Terms of Service</a> apply.
      </p>
    </b-card>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, email } from 'vuelidate/lib/validators'
import { regex } from '~/assets/config'
import Loading from '~/components/PageLoading.vue'
const validPassword = val => regex.password.test(val)
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  name: 'SignUp',
  components: {
    Loading
  },
  mixins: [validationMixin],
  data() {
    return {
      loading: true,
      sendEmailMode: true,
      emailform: {
        email: ''
      },
      passwordform: {
        password: ''
      }
    }
  },
  // @ts-ignore
  validations: {
    emailform: {
      email: {
        required,
        email
      }
    },
    passwordform: {
      password: {
        required,
        validPassword
      }
    }
  },
  // @ts-ignore
  head() {
    const title = 'Reset'
    const description = 'reset your account password'
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
  mounted() {
    if (
      this.$route.query &&
      this.$route.query.reset &&
      this.$route.query.token
    ) {
      this.sendEmailMode = false
    }
    this.loading = false
  },
  methods: {
    resetPassword(evt) {
      evt.preventDefault()
      this.$axios
        .post('/reset', {
          token: this.$route.query.token,
          password: this.passwordform.password
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              let message = 'password reset successfully'
              if (res.data.message) {
                message = res.data.message
              }
              this.$toasted.global.success({
                message: message
              })
              this.$router.push({
                path: '/login'
              })
            } else {
              this.$toasted.global.error({
                message: 'could not get data'
              })
            }
          } else if (res.data && res.data.message) {
            this.$toasted.global.error({
              message: res.data.message
            })
          } else {
            this.$toasted.global.error({
              message: `status code of ${res.status}`
            })
          }
        })
        .catch(err => {
          /* eslint-disable */
          let message = `got error: ${err}`
          if (err.response && err.response.data) {
            message = err.response.data.message
          }
          this.$toasted.global.error({
            message: message
          })
        })
    },
    sendResetEmail(evt) {
      evt.preventDefault()
      this.$recaptcha('login').then(recaptchatoken => {
        this.$axios
          .put('/sendResetEmail', {
            email: this.emailform.email,
            recaptcha: recaptchatoken
          })
          .then(res => {
            if (res.status === 200) {
              if (res.data) {
                let message = 'email sent'
                if (res.data.message) {
                  message = res.data.message
                }
                this.$toasted.global.success({
                  message: message
                })
                this.emailform = {
                  email: ''
                }
              } else {
                this.$toasted.global.error({
                  message: 'could not get data'
                })
              }
            } else if (res.data && res.data.message) {
              this.$toasted.global.error({
                message: res.data.message
              })
            } else {
              this.$toasted.global.error({
                message: `status code of ${res.status}`
              })
            }
          })
          .catch(err => {
            /* eslint-disable */
            let message = `got error: ${err}`
            if (err.response && err.response.data) {
              message = err.response.data.message
            }
            this.$toasted.global.error({
              message: message
            })
          })
      }).catch(err => {
        this.$toasted.global.error({
          message: `got error with recaptcha ${err}`
        })
      })
    }
  }
})
</script>

<style lang="scss"></style>
