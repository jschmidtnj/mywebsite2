<template>
  <b-container class="mt-4">
    <b-card title="Sign up" footer-tag="footer">
      <b-form @submit="signup" @reset="reset">
        <b-form-group
          id="email-address-group"
          label="Email address:"
          label-for="email-address"
          description="Your email is safe with us"
        >
          <b-form-input
            id="email-address"
            v-model="form.email"
            :state="!$v.form.email.$invalid"
            type="text"
            autocomplete="off"
            placeholder="Enter email"
            aria-describedby="emailfeedback"
          ></b-form-input>
          <b-form-invalid-feedback
            id="emailfeedback"
            :state="!$v.form.email.$invalid"
          >
            <div v-if="!$v.form.email.required">email is required</div>
            <div v-else-if="!$v.form.email.email">
              email is invalid
            </div>
          </b-form-invalid-feedback>
        </b-form-group>
        <b-form-group
          id="password-group"
          label="Password:"
          label-for="password"
          description="Password must be at least 8 characters long, with a number, capital letter and special character"
        >
          <b-form-input
            id="password"
            v-model="form.password"
            :state="!$v.form.password.$invalid"
            type="password"
            autocomplete="off"
            placeholder="Enter password"
            aria-describedby="passwordfeedback"
          ></b-form-input>
          <b-form-invalid-feedback
            id="passwordfeedback"
            :state="!$v.form.password.$invalid"
          >
            <div v-if="!$v.form.password.required">password is required</div>
            <div v-else-if="!$v.form.password.validPassword">
              password is invalid
            </div>
          </b-form-invalid-feedback>
        </b-form-group>
        <b-button
          :disabled="$v.form.$invalid"
          variant="primary"
          type="submit"
          class="mt-4"
          >Submit</b-button
        >
      </b-form>
      <p slot="footer">
        By clicking submit you aggree to the
        <nuxt-link to="/privacy">privacy policy</nuxt-link>. This site is
        protected by reCAPTCHA and the Google
        <a href="https://policies.google.com/privacy">Privacy Policy</a> and
        <a href="https://policies.google.com/terms">Terms of Service</a> apply.
      </p>
    </b-card>
  </b-container>
</template>

<script lang="js">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, email } from 'vuelidate/lib/validators'
import { regex } from '~/assets/config'
const validPassword = val => regex.password.test(val)
// @ts-ignore
const seo = JSON.parse(process.env.seoconfig)
export default Vue.extend({
  name: 'SignUp',
  mixins: [validationMixin],
  data() {
    return {
      form: {
        email: '',
        password: ''
      }
    }
  },
  // @ts-ignore
  validations: {
    form: {
      email: {
        required,
        email
      },
      password: {
        required,
        validPassword
      }
    }
  },
  methods: {
    reset(evt) {
      evt.preventDefault()
      this.form = {
        email: '',
        password: ''
      }
    },
    signup(evt) {
      evt.preventDefault()
      this.$recaptcha('login')
        .then(recaptchatoken => {
          this.$axios
            .post('/register', {
              email: this.form.email,
              password: this.form.password,
              recaptcha: recaptchatoken
            })
            .then(res => {
              if (res.status === 200) {
                if (res.data) {
                  let message =
                    'finished signing up. please check email for confirmation'
                  if (res.data.message) {
                    message = res.data.message
                  }
                  this.$toasted.global.success({
                    message
                  })
                  this.reset(evt)
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
              let message = `got error: ${err}`
              if (err.response && err.response.data) {
                message = err.response.data.message
              }
              this.$toasted.global.error({
                message
              })
            })
        })
        .catch(err => {
          this.$toasted.global.error({
            message: `got error with recaptcha ${err}`
          })
        })
    }
  },
  // @ts-ignore
  head() {
    const title = 'Sign Up'
    const description = 'sign up for an account'
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

<style lang="scss"></style>
