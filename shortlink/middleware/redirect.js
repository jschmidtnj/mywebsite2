/* eslint-disable */
export default ({ redirect, route, env }) => {
  const path = route.path.substr(1)
  redirect(`${env.apiurl}/shortlink?id=${path}`)
}
