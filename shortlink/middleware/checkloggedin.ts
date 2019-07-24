/* eslint-disable */

export default ({ store, redirect }) => {
  return new Promise((resolve, reject) => {
    if (!store.state.auth) {
      resolve()
    } else {
      store.dispatch('auth/checkLoggedIn').then(loggedin => {
        if (!loggedin) {
          resolve()
        } else if (!store.state.auth.user) {
          store.dispatch('auth/getUser').then(res => {
            resolve()
          }).catch(err => {
            resolve()
          })
        } else {
          resolve()
        }
      }).catch(err => {
        resolve()
      })
    }
  })
}
