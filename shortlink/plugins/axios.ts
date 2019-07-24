export default ({ $axios, store }) => {
  if (store.state.auth && store.state.auth.token) {
    $axios.setToken(store.state.auth.token)
  }
}
