/* eslint-disable */

export default ({ redirect }) => {
  const mainurl = process.env.mainurl ? process.env.mainurl : ''
  redirect(`${mainurl}/signup`)
}
