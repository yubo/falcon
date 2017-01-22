import { MUT_NOTIFICATION_SET, MUT_NOTIFICATION_CLOSE } from '../mutation-types'

// initial state
// shape: [{ id, quantity }]
const state = {
  show: false,
  msg: '',
  level: 'INFO'
}

// getters
const getters = {
}

// actions
const actions = {
  notification_set ({ commit }, args) {
    commit(MUT_NOTIFICATION_SET, args)
    if (args.timeout > 0) {
      setTimeout(() => {
        commit(MUT_NOTIFICATION_CLOSE)
      }, args.timeout)
    }
  },
  notification_close ({ commit }) {
    commit(MUT_NOTIFICATION_CLOSE)
  }
}

// mutations
const mutations = {
  [MUT_NOTIFICATION_SET] (state, { level, msg }) {
    state.level = level
    state.msg = msg
    state.show = true
  },
  [MUT_NOTIFICATION_CLOSE] (state) {
    state.show = false
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
