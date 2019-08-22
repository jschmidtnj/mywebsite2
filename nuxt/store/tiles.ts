/**
 * tile carousel data for persistance through page changes
 */

export const state = () => ({
  projects: [],
  blogs: [],
  blogcount: 0,
  projectcount: 0,
  blogindex: null,
  projectindex: null,
  perpage: 8,
  sortBy: 'date',
  sortDesc: true,
  tags: [],
  categories: []
})

export const getters = {
  projects: state => state.projects,
  blogs: state => state.blogs,
  perpage: state => state.perpage,
  sortBy: state => state.sortBy,
  sortDesc: state => state.sortDesc
}

export const mutations = {
  setPosts(state, payload) {
    if (payload.type === 'blog') {
      state.blogs = payload.posts
    } else {
      state.projects = payload.posts
    }
  },
  addPosts(state, payload) {
    if (payload.type === 'blog') {
      state.blogs[payload.index] = payload.posts
    } else {
      state.projects[payload.index] = payload.posts
    }
  },
  setCount(state, payload) {
    if (payload.type === 'blog') {
      state.blogcount = payload.count
    } else {
      state.projectcount = payload.count
    }
  },
  setIndex(state, payload) {
    if (payload.type === 'blog') {
      state.blogindex = payload.index
    } else {
      state.projectindex = payload.index
    }
  }
}

export const actions = {
  async updateCount({ state, commit, rootState }, payload) {
    return new Promise((resolve, reject) => {
      this.$axios
        .get('/countPosts', {
          params: {
            searchterm: '',
            type: payload.type,
            tags: state.tags.join(',tags='),
            categories: state.categories.join(',categories='),
            cache: !(
              rootState.auth.user &&
              rootState.auth.user.type === 'admin'
            )
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.count !== null) {
                commit('setCount', {
                  type: payload.type,
                  count: res.data.count
                })
                resolve(`count updated successfully to ${res.data.count}`)
              } else {
                reject('could not find count data')
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
  },
  initializePosts({ state, commit }, payload) {
    if (payload.type === 'blog') {
      const numPages = Math.ceil(state.blogcount / state.perpage)
      commit('setPosts', {
        type: 'blog',
        posts: Array.apply(null, Array(numPages))
      })
    } else {
      const numPages = Math.ceil(state.projectcount / state.perpage)
      commit('setPosts', {
        type: 'project',
        posts: Array.apply(null, Array(numPages))
      })
    }
  },
  async addPosts({ state, commit, rootState }, payload) {
    return new Promise((resolve, reject) => {
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${encodeURIComponent(
              payload.type
            )}",perpage:${encodeURIComponent(
              state.perpage
            )},page:${
              payload.page
              },searchterm:"",sort:"${encodeURIComponent(
                state.sortBy
              )}",ascending:${
              !state.sortDesc
              },tags:${
              JSON.stringify(state.tags)
              },categories:${
              JSON.stringify(state.categories)
              },cache:${(!(
                rootState.auth.user &&
                rootState.auth.user.type === 'admin'
              )).toString()}){tileimage id title caption color}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                commit('addPosts', {
                  type: payload.type,
                  index: payload.page,
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
