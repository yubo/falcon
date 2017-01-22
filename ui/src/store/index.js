import Vue from 'vue'
import Vuex from 'vuex'
import login from './modules/login'
import notification from './modules/notification'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    tree: {
      show: true
    }
  },
  // actions,
  // getters,
  modules: {
    login,
    notification
  }
})
