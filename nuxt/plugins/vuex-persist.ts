import createPersistedState from 'vuex-persistedstate'

export default ({ store }) => {
  // @ts-ignore
  window.onNuxtReady(() => {
    createPersistedState({
      key: 'yourkey',
      reducer: persistedState => {
        const stateFilter = Object.assign({}, persistedState)
        stateFilter.auth = {}
        return stateFilter
      }
    })(store)
  })
}
