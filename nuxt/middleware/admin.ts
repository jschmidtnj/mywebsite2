export default ({ store, redirect }) => {
  console.log(`auth: ${JSON.stringify(store.state.auth)}`)
  if (store.state.auth && store.state.auth.loggedIn && 
      store.state.auth.user.emailverified && store.state.auth.user.type === 'admin') {
    console.log('signed in')
  } else {
    console.log('not signed in')
    redirect('/login')
  }
}