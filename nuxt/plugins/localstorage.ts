import createPersistedState from 'vuex-persistedstate'

export default ({ store }) => {
  // @ts-ignore
  window.onNuxtReady(() => {
    createPersistedState({
      key: 'yourkey'
    })(store)
  })
}
