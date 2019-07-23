import * as jwt from 'jsonwebtoken'

/**
 * authentication
 */

export const state = () => ({
  token: null,
  user: null,
  loggedIn: false
})

export const getters = {
  token: state => state.token,
  user: state => state.user
}

export const mutations = {
  setToken(state, payload) {
    state.token = payload
  },
  setUser(state, payload) {
    state.user = payload
  },
  setLoggedIn(state, payload) {
    state.loggedIn = payload
  }
}

export const actions = {
  checkLoggedIn({state, commit}) {
    let res = true
    try {
      const { exp } = jwt.decode(state.token)
      if (Date.now() >= exp * 1000) {
        res = false
      }
    } catch (err) {
      res = false
    }
    commit('setLoggedIn', res)
    return res
  },
  setToken({ state, commit, app: { $axios } }, token) {
    commit('setToken', token)
    $axios.setToken(state.token)
    return 'got token'
  },
  async getUser({ state, commit, app: { $axios } }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        reject('no token found for user')
      } else {
        $axios.get('/graphql', {
          params: {
            query: '{account{id email type emailverified shortlinks}}'
          }
        }).then(res => {
          if (res.status === 200) {
            if (res.data) {
              if (res.data.data && res.data.data.account) {
                commit('setUser', res.data.data.account)
                resolve('found user account data')
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
        }).catch(err => {
          reject(err)
        })
      }
    })
  }
}