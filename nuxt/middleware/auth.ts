export default ({ store, redirect }) => {
  if (store.state.auth && store.state.auth.loggedIn) {
    console.log('signed in')
  } else {
    console.log('not signed in')
    redirect('/login')
  }
}