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
  updateUser(state, payload) {
    let schema = state.user
    const pList = payload.path.split('.')
    const len = pList.length
    for (let i = 0; i < len - 1; i++) {
      const elem = pList[i]
      if (!schema[elem]) schema[elem] = {}
      schema = schema[elem]
    }
    if (payload.value && payload.value === 'DELETE') {
      delete schema[pList[len - 1]]
    } else {
      schema[pList[len - 1]] = payload.value
    }
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
  checkLoggedIn({ state, commit }) {
    let res = true
    try {
      const decoded = jwt.decode(state.token, {
        complete: true
      })
      if (Date.now() >= decoded.exp * 1000) {
        res = false
      }
    } catch (err) {
      res = false
    }
    commit('setLoggedIn', res)
    return res
  },
  async getUser({ state, commit }) {
    return new Promise((resolve, reject) => {
      if (!state.token) {
        reject(new Error('no token found for user'))
      } else {
        this.$axios
          .get('/graphql', {
            params: {
              query: '{account{id email type emailverified shortlinks}}'
            }
          })
          .then(res => {
            if (res.status === 200) {
              if (res.data) {
                if (res.data.data && res.data.data.account) {
                  commit('setUser', res.data.data.account)
                  resolve('found user account data')
                } else if (res.data.errors) {
                  reject(
                    new Error(
                      `found errors: ${JSON.stringify(res.data.errors)}`
                    )
                  )
                } else {
                  reject(new Error('could not find data or errors'))
                }
              } else {
                reject(new Error('could not get data'))
              }
            } else {
              reject(new Error(`status code of ${res.status}`))
            }
          })
          .catch(err => {
            reject(err)
          })
      }
    })
  }
}
