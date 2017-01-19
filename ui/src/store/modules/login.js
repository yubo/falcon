import * as types from '../mutation-types'
import { fetch } from '../utils'

// initial state
// shape: [{ id, quantity }]
const state = {
  notification: '',
  status: false
}

// getters
const getters = {
}

// actions
const actions = {
  login ({ commit, state }, args) {
    fetch({
      method: 'post',
      url: 'auth/login',
      params: {
        username: args.username,
        password: args.password,
        method: args.method
      }
    })
    .then((res) => {
      commit(types.LOGIN_SUCCESS, {
        data: res.data,
        router: args.router
      })
    })
    .catch((err) => {
      commit(types.LOGIN_FAIL, {
        err
      })
    })
  }
}

// mutations
const mutations = {
  [types.LOGIN_SUCCESS] (state, { data, router }) {
    window.Cookies.set('name', data.name)
    window.Cookies.set('sig', data.sig)
    state.notification = 'Log in Success!'
    state.status = true
    setTimeout(() => {
      router.push('/graph')
    }, 2000)
  },
  [types.LOGIN_FAIL] (state, { err }) {
    state.notification = 'Username or Password error'
    state.status = false
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
