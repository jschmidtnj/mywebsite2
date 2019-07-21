/* eslint-disable */

export default ({ store, redirect, query }) => {
  if (query && query.token) {
    store.dispatch('auth/setToken', query.token).then(res => {
      console.log(res)
    }).catch(err => {
      console.log(err)
      redirect('/')
    })
  } else {
    console.log('no token found')
    redirect('/')
  }
}