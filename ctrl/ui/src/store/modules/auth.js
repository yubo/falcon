import { fetch, Msg } from 'src/utils'

const state = {
  user: null,
  callback: '',
  admin: false,
  reader: false,
  operator: false,
  login: false
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
      Msg.success('logout success')
      console.log(args)
      window.location.href = '/'
    }).catch((err) => {
      Msg.error('logout failed', err)
    })
  },
  info ({ commit, state }) {
    return new Promise((resolve, reject) => {
      fetch({
        url: 'auth/info',
        method: 'get'
      }).then((res) => {
        if (res.data.user) {
          commit('m_login_success', res.data)
          Msg.success('welecom ' + res.data.user.name)
        }
        resolve(res)
      }).catch((err) => {
        reject(err)
      })
    })
  },
  login ({ commit, state }, args = {}) {
    return new Promise((resolve, reject) => {
      fetch({
        url: 'auth/login',
        method: 'post',
        params: {
          username: args.username,
          password: args.password,
          method: args.method
        }
      }).then((res) => {
        commit('m_login_success', res.data)
        resolve(res)
      }).catch((err) => {
        commit('m_login_fail')
        Msg.error('login fail', err)
        reject(err)
      })
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
