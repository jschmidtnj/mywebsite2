<template>
  <b-card title="Reset" footer-tag="footer">
    <b-form @submit="reset" @reset="clear">
      <b-form-group
        id="email-address-group"
        label="Email address:"
        label-for="email-address"
        description="Please enter your account email address"
      >
        <b-form-input
          id="email-address"
          v-model="form.email"
          type="text"
          autocomplete="off"
          :state="!$v.form.email.$invalid"
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
      <a href="/privacy">privacy policy</a>.
    </p>
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { validationMixin } from 'vuelidate'
import { required, email } from 'vuelidate/lib/validators'
export default Vue.extend({
  name: 'SignUp',
  mixins: [validationMixin],
  data() {
    return {
      form: {
        email: ''
      }
    }
  },
  // @ts-ignore
  validations: {
    form: {
      email: {
        required,
        email
      }
    }
  },
  methods: {
    clear(evt) {
      evt.preventDefault()
      this.form = {
        email: '',
        password: ''
      }
    },
    reset(evt) {
      evt.preventDefault()
      this.$axios
        .put('/sendResetEmail', {
          password: this.form.password
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
              this.reset()
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
    }
  }
})
</script>

<style lang="scss"></style>
