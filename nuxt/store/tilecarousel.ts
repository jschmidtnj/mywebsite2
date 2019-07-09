/**
 * tile carousel data for persistance through page changes
 */

export const state = () => ({
  projectpage: 0,
  blogpage: 0,
  projects: [],
  blogs: [],
  blogcount: 0,
  projectcount: 0,
  perpage: 8,
  sortBy: 'title',
  sortDesc: true,
  tags: [],
  categories: []
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
  },
  setCount(state, payload) {
    if (payload.type === 'blog') {
      state.blogcount = payload.count
    } else {
      state.projectcount = payload.count
    }
  }
}

export const actions = {
  async updateCount({ state, commit }, payload) {
    console.log(`got state: `)
    console.log(state)
    return new Promise((resolve, reject) => {
      this.$axios
        .get('/countPosts', {
          params: {
            searchterm: '',
            type: payload.type,
            tags: state.tags.join(',tags='),
            categories: state.categories.join(',categories=')
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.count !== null) {
                commit('setCount', {
                  type: payload.type,
                  posts: res.data.count
                })
                resolve('count updated successfully')
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
  async updatePosts({ state, commit }, payload) {
    console.log(`got state: `)
    console.log(state)
    return new Promise((resolve, reject) => {
      this.$axios
        .get('/graphql', {
          params: {
            query: `{posts(type:"${encodeURIComponent(
              payload.type
            )}",perpage:${encodeURIComponent(
              state.perpage
            )},page:${
              payload.type === 'blog' ? state.blogpage : state.projectpage
            },searchterm:"",sort:"${encodeURIComponent(
              state.sortBy
            )}",ascending:${
              !state.sortDesc
            },tags:${
              JSON.stringify(state.tags)
            },categories:${
              JSON.stringify(state.categories)
            },cache:true){tileimage id title caption}}`
          }
        })
        .then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.posts) {
                commit('setPosts', {
                  type: payload.type,
                  posts: res.data.data.posts.filter(post => post.tileimage.length > 0)
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
