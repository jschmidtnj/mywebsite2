/* eslint-disable */

export default ({ store, redirect }) => {
  return new Promise((resolve, reject) => {
    if (!store.state.auth) {
      redirect('/login')
    } else {
      store.dispatch('auth/checkLoggedIn').then(loggedin => {
        if (!loggedin) {
          redirect('/login')
        } else if (!store.state.auth.user) {
          store.dispatch('auth/getUser').then(res => {
            resolve()
          }).catch(err => {
            redirect('/login')
          })
        } else {
          resolve()
        }
      }).catch(err => {
        redirect('/login')
      })
    }
  })
}
