export default ({ store, redirect }) => {
  console.log(`auth: ${JSON.stringify(store.state.auth)}`)
  if (store.state.auth && store.state.auth.loggedIn) {
    console.log('signed in')
    redirect('/account')
  }
}