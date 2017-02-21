import { fetch } from 'src/utils'
import { Message } from 'element-ui'

const state = {
  user: null,
  callback: '',
  admin: false,
  reader: false,
  operator: false,
  login: false,
  loading: false
}

const getters = {
}

const actions = {
  logout ({ commit }) {
    fetch({
      method: 'get',
      url: 'auth/logout'
    }).then((res) => {
      commit('m_logout')
      Message.success('logout success')
    })
    .catch((err) => {
      Message.error(err.response.data)
    })
  },
  info ({ commit, state }) {
    fetch({
      url: 'auth/info',
      method: 'get'
    }).then((res) => {
      if (res.data.user) {
        commit('m_login_success', res.data)
        Message.success('welecom ' + res.data.user.name)
      }
    })
  },
  login ({ commit, state }, args = {}) {
    commit('m_set_loading', true)
    fetch({
      url: 'auth/login',
      method: 'post',
      params: args
    }).then((res) => {
      commit('m_login_success', res.data)
      commit('m_set_loading', false)
      Message.success('login success, hi ' + res.data.user.name)
    }).catch((err) => {
      commit('m_login_fail')
      commit('m_set_loading', false)
      Message.error('login fail ' + err)
    })
  }
}

// mutations
const mutations = {
  'm_set_user' (state, user) {
    state.user = user
    state.username = state.user.name
  },
  'm_logout' (state) {
    state.login = false
    state.reader = false
    state.operator = false
    state.admin = false
    window.Cookies.remove('username')
  },
  'm_set_loading' (state, loading) {
    state.loading = loading
  },
  'm_login_success' (state, obj) {
    state.login = true
    state.user = obj.user
    state.username = state.user.name
    state.reader = obj.reader
    state.operator = obj.operator
    state.admin = obj.admin
    window.Cookies.set('username', obj.name, {expires: 1})
  },
  'm_login_fail' (state) {
    state.login = false
    window.Cookies.remove('username')
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
