<template>
  <b-card>
    <b-table striped hover :items="shortlinks" :fields="fields"></b-table>
  </b-card>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  name: 'Account',
  // @ts-ignore
  layout: 'secure',
  data: function() {
    return {
      fields: {
        id: {
          key: 'id',
          label: 'Id',
          sortable: true
        },
        link: {
          key: 'link',
          label: 'Link',
          sortable: true
        }
      },
      shortlinks: []
    }
  },
  mounted() {
    if (
      this.$store.state.auth &&
      this.$store.state.auth.user &&
      this.$store.state.auth.user.shortlinks
    ) {
      this.$axios
        .get('/graphql', {
          params: {
            query: `{shortlinks(linkids:${encodeURIComponent(
              this.$store.state.auth.user.shortlinks
            )}){link id}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.shortlinks) {
                this.shortlinks = res.data.data.shortlinks
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
    } else {
      this.$toasted.global.error({
        message: `no shortlink ids found`
      })
    }
  }
})
</script>

<style lang="scss"></style>
