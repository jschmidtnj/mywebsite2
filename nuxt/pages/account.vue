<template>
  <b-card title="Account Page">
    <p>user: {{ $auth.user }}</p>
    <b-btn @click="$auth.logout()">Logout</b-btn>
    <b-btn @click="deleteAccount">Delete</b-btn>
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  // @ts-ignore
  layout: 'secure',
  // @ts-ignore
  head() {
    return {
      title: 'Account'
    }
  },
  methods: {
    deleteAccount(evt) {
      evt.preventDefault()
      this.$axios
        .delete('/graphql', {
          params: {
            query: `mutation{deleteAccount(){id}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.deleteAccount) {
                this.$toasted.global.success({
                  message: 'account deleted'
                })
                this.$router.push({
                  path: '/login'
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
