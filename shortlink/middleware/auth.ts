/* eslint-disable */

export default ({ store, redirect }) => {
  return new Promise((resolve, reject) => {
    if (!store.state.auth) {
      console.log('could not get auth data')
      redirect('/')
    } else {
      console.log(`auth: ${JSON.stringify(store.state.auth)}`)
      store.dispatch('auth/checkLoggedIn').then(loggedin => {
        console.log(loggedin)
        if (!loggedin) {
          console.log('not logged in')
          redirect('/')
        } else if (!store.state.auth.user) {
          store.dispatch('auth/getUser').then(res => {
            console.log(res)
            resolve()
          }).catch(err => {
            console.log(err)
            redirect('/')
          })
        } else {
          resolve()
        }
      }).catch(err => {
        console.log(err)
        redirect('/')
      })
    }
  })
}
