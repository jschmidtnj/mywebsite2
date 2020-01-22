/* eslint-disable */

export default ({ redirect }) => {
  // @ts-ignore
  const mainurl = process.env.mainurl
  // @ts-ignore
  const currenturl = JSON.parse(process.env.seoconfig).url
  const redirecturl = encodeURIComponent(`${currenturl}/callback`)
  redirect(`${mainurl}/login?redirect_uri=${redirecturl}`)
}
