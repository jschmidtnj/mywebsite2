/* eslint-disable */

export default ({ redirect }) => {
  const mainurl = process.env.mainurl ? process.env.mainurl : ''
  const currenturl = process.env.seoconfig
    ? JSON.parse(process.env.seoconfig).url
    : ''
  const redirecturl = encodeURIComponent(`${currenturl}/callback`)
  redirect(`${mainurl}/login?redirect_uri=${redirecturl}`)
}
