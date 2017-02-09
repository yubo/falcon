import { vfetch, fetch } from 'src/utils'
import { Message } from 'element-ui'
const { _ } = window

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
    var opts = {
      commit,
      mutation: 'm_login',
      args: args,
      method: 'post',
      url: 'auth/login'
    }
    if (args.username) {
      opts.params = {
        username: args.username,
        password: args.password,
        method: args.method
      }
    }
    vfetch(opts)
  }
}

// mutations
const mutations = {
  'm_set_user' (state, user) {
    state.user = user
  },
  'm_set_callback' (state, cb) {
    state.callback = cb
  },
  'm_logout' (state) {
    state.login = false
    window.Cookies.remove('username')
  },
  'm_login.start' (state) {
    // console.log('login.start')
    state.loading = true
  },
  'm_login.success' (state, args = {}) {
    // console.log('login.success', args)
    state.login = true
    state.user = args.res.data
    window.Cookies.set('username', args.res.data.name, {expires: 1})
    Message.success('login success')
    if (args.router) {
      if (state.callback) {
        args.router.push(state.callback)
        state.callback = ''
      } else {
        args.router.push('/meta')
      }
    }
  },
  'm_login.fail' (state, args) {
    state.login = false
    window.Cookies.remove('username')
    Message.error(args.err)
    if (args.router && (!_.startsWith(args.router.path, '/login'))) {
      state.callback = args.router.fullPath
      args.router.push({path: '/login'})
    }
  },
  'm_login.end' (state) {
    state.loading = false
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
