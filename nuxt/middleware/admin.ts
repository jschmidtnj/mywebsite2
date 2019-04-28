export default ({ store, redirect }) => {
  // console.log(`auth: ${JSON.stringify(store.state.auth)}`)
  // use custom strategy with custom token that I can
  // read from the json
  if (store.state.auth && store.state.auth.loggedIn) {
    console.log('signed in')
  } else {
    console.log('not signed in')
    redirect('/login')
  }
}
