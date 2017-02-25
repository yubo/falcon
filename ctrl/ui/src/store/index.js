import Vue from 'vue'
import Vuex from 'vuex'
import auth from './modules/auth'
import rel from './modules/rel'
import { fetch, Msg } from 'src/utils'

Vue.use(Vuex)

const modules = { auth, rel }

const state = {
  config: null
}

const getters = {
  schema: state => {
    if (state.config[2] && state.config[2]['tagschema']) {
      return state.config[2]['tagschema']
    }
    if (state.config[1] && state.config[1]['tagschema']) {
      return state.config[1]['tagschema']
    }
    if (state.config[0] && state.config[0]['tagschema']) {
      return state.config[0]['tagschema']
    }
    return ''
  }
}

const actions = {
  load_config ({state, commit}) {
    fetch({
      method: 'get',
      url: 'settings/config/ctrl'
    }).then((res) => {
      commit('m_set_config', res.data)
    }).catch((err) => {
      Msg.error('get config failed', err)
    })
  }
}

const mutations = {
  'm_set_config' (state, config) {
    state.config = config
  }
}

const store = new Vuex.Store({
  modules,
  state,
  getters,
  actions,
  mutations
})

export default store
