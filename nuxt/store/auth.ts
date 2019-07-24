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
    this.$axios.setToken(state.token)
  },
  setUser(state, payload) {
    state.user = payload
  },
  setLoggedIn(state, payload) {
    state.loggedIn = payload
  },
  logout(state) {
    state.token = null
    state.user = null
    state.loggedIn = false
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
  async loginLocal({ commit }, payload) {
    return new Promise((resolve, reject) => {
      this.$axios.put('/loginEmailPassword', {
        email: payload.email,
        password: payload.password,
        recaptcha: payload.recaptcha
      }).then(res => {
        if (res.status === 200) {
          if (res.data) {
            if (res.data.token) {
              resolve(res.data.token)
            } else {
              reject('could not find token data')
            }
          } else {
            reject('could not get data')
          }
        } else {
          reject(`status code of ${res.status}`)
        }
      }).catch(err => {
        let message = `got error: ${err}`
        if (err.response && err.response.data) {
          message = err.response.data.message
        }
        reject(message)
      })
    })
  },
  async getUser({ state, commit }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        reject('no token found for user')
      } else {
        this.$axios.get('/graphql', {
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