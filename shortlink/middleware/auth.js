/* eslint-disable */

export default ({ store, redirect }) => {
  return new Promise((resolve, reject) => {
    if (!store.state.auth) {
      redirect('/')
    } else {
      store.dispatch('auth/checkLoggedIn').then(loggedin => {
        if (!loggedin) {
          redirect('/')
        } else if (!store.state.auth.user) {
          store.dispatch('auth/getUser').then(res => {
            resolve()
          }).catch(err => {
            redirect('/')
          })
        } else {
          resolve()
        }
      }).catch(err => {
        redirect('/')
      })
    }
  })
}
