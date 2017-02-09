import Vue from 'vue'
import Vuex from 'vuex'
import login from './modules/login'
import rel from './modules/rel'

Vue.use(Vuex)

const modules = { login, rel }

const state = {
  config: null
}

const getters = {
}

const actions = {
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
