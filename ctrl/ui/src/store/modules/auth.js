import { fetch } from 'src/utils'
import { Message } from 'element-ui'

const state = {
  user: null,
  callback: '',
  login: false,
  loading: false
}

const getters = {
}

const actions = {
  logout ({ commit }, args = {}) {
    fetch({
      method: 'get',
      url: 'auth/logout'
    }).then((res) => {
      commit('m_logout')
      Message.success('logout success')
      if (args.router) {
        args.router.push('/login')
      }
    })
    .catch((err) => {
      Message.error(err.response.data)
    })
  },
  login ({ commit, state }, args = {}) {
    commit('m_set_loading', true)
    fetch({
      url: 'auth/login',
      method: 'post',
      params: args
    }).then((res) => {
      commit('m_login_success', res)
      commit('m_set_loading', false)
      Message.success('login success, hi ' + res.data.name)
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
    window.Cookies.remove('username')
  },
  'm_set_loading' (state, loading) {
    state.loading = loading
  },
  'm_login_success' (state, res) {
    state.login = true
    state.user = res.data
    state.username = state.user.name
    window.Cookies.set('username', res.data.name, {expires: 1})
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
