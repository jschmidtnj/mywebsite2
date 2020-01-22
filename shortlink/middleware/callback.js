/* eslint-disable */

export default ({ store, redirect, query }) => {
  return new Promise((resolve, reject) => {
    if (query && query.token) {
      store.commit('auth/setToken', query.token)
      store.dispatch('auth/getUser').then(res => {
        redirect('/account')
      }).catch(err => {
        redirect('/')
      })
    } else {
      redirect('/')
    }
  })
}