export default ({ $axios, store }) => {
  $axios.onRequest(config => {
    if (store.state.auth.strategy) {
      console.log(`got strategy ${store.state.auth.strategy}`)
      config.headers.common['Strategy'] = store.state.auth.strategy
    }
  })
}