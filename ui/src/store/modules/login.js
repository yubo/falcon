import { fetch } from '../utils'
import {
  MUT_LOGIN,
  MUT_LOGOUT,
  MUT_NOTIFICATION_SET
} from '../mutation-types'

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
  logout ({ commit, state }) {
    fetch({
      method: 'get',
      url: 'auth/logout'
    }).then((res) => {
      commit(MUT_LOGOUT)
      commit(MUT_NOTIFICATION_SET, { level: 'SUCCESS', msg: 'logout Success!' })
    })
    .catch((err) => {
      commit(MUT_NOTIFICATION_SET, { level: 'ERROR', msg: err.response.data })
    })
  },
  login_quiet ({ commit, state }) {
    fetch({
      method: 'post',
      url: 'auth/login'
    }).then(() => {
      commit(MUT_LOGIN)
    })
  },
  login ({ commit, state }, args) {
    fetch({
      method: 'post',
      url: 'auth/login',
      params: {
        username: args.username,
        password: args.password,
        method: args.method
      }
    }).then((res) => {
      commit(MUT_LOGIN)
      commit(MUT_NOTIFICATION_SET, { level: 'SUCCESS', msg: 'login Success!' })
      setTimeout(() => { args.router.push('/meta/tag') }, 2000)
    }).catch((err) => {
      commit(MUT_NOTIFICATION_SET, { level: 'ERROR', msg: err.response.data })
    })
  }
}

// mutations
const mutations = {
  [MUT_LOGOUT] (state) { state.status = false },
  [MUT_LOGIN] (state) { state.status = true }
}

export default {
  state,
  getters,
  actions,
  mutations
}
