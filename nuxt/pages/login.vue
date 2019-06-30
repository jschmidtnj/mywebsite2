<template>
  <div>
    <p>login</p>
    <b-btn block @click="logingoogle">Login with Google</b-btn>
    <b-form @submit="loginlocal">
      <b-form-group
        label="Email address:"
        label-for="emailinput"
        description="We'll never share your email with anyone else"
      >
        <b-form-input
          id="emailinput"
          v-model="email"
          type="email"
          required
          placeholder="Enter email"
        ></b-form-input>
      </b-form-group>
      <b-form-group
        label="Password"
        label-for="passwordinput"
        description="Please enter valid password"
      >
        <b-form-input
          id="passwordinput"
          v-model="password"
          type="password"
          required
          placeholder
        ></b-form-input>
      </b-form-group>
      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'Login',
  // @ts-ignore
  middleware: 'loginredirect',
  data() {
    return {
      email: '',
      password: ''
    }
  },
  // @ts-ignore
  head() {
    return {
      title: 'Login'
    }
  },
  methods: {
    logingoogle() {
      /* eslint-disable */
      this.$auth.loginWith('google').catch(e => {
        console.log(e)
      })
    },
    loginlocal(evt) {
      evt.preventDefault()
      this.$auth.loginWith('local', {
        data: {
          email: this.email,
          password: this.password
        }
      }).then(() => {
        console.log('success')
      }).catch(err => {
        console.log(`got error logging in local: ${err}`)
      })
    }
  }
})
</script>

<style lang="scss"></style>