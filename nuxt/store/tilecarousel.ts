/**
 * tile carousel data for persistance through page changes
 */

export const state = () => ({
  projectpage: 0,
  blogpage: 0,
  projects: [],
  blogs: [],
  perpage: 5,
  sortBy: 'title',
  sortDesc: true
})

export const getters = {
  projectpage: state => state.projectpage,
  blogpage: state => state.blogpage,
  projects: state => state.projects,
  blogs: state => state.blogs,
  perpage: state => state.perpage,
  sortBy: state => state.sortBy,
  sortDesc: state => state.sortDesc
}

export const mutations = {
  setPage(state, payload) {
    if (payload.type === 'blog') {
      state.blogpage = payload.page
    } else {
      state.projectpage = payload.page
    }
  },
  setPosts(state, payload) {
    if (payload.type === 'blog') {
      state.blogs = payload.posts
    } else {
      state.projects = payload.posts
    }
  }
}

export const actions = {
  async updateCarousel({ state, commit }, payload) {
    return new Promise((resolve, reject) => {
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${encodeURIComponent(payload.type)}",perpage:${
              state.perPage
            },page:${
              payload.type === 'blog' ? state.blogpage : state.projectpage
            },searchterm:"",sort:"${encodeURIComponent(
              state.sortBy
            )}",ascending:${
              !state.sortDesc
            },tags:${JSON.stringify([])},categories:${JSON.stringify(
              []
            )},cache:true){tile id title caption}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                commit('setPosts', {
                  type: payload.type,
                  posts: res.data.data.posts
                })
                resolve('found posts')
              } else if (res.data.errors) {
                reject(`found errors: ${JSON.stringify(res.data.errors)}`)
              } else {
                reject('could not find data or errors')
              }
            } else {
              reject('could not get data')
            }
          } else {
            reject(`status code of ${res.status}`)
          }
        })
        .catch(err => {
          let message = `got error: ${err}`
          if (err.response && err.response.data) {
            message = err.response.data.message
          }
          reject(message)
        })
    })
  }
}
