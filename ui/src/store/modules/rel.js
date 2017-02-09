import { fetch } from 'src/utils'
import { Message } from 'element-ui'

const state = {
  loading: false,
  loaded: false,
  curTag: {name: '', id: 1},
  tree: []
}

const getters = {
}

const actions = {
}

const mutations = {
  'm_cur_tag' (state, val) {
    state.curTag = val
  },
  'm_load_tag' (state, router) {
    state.loading = true
    fetch({
      router: router,
      method: 'get',
      url: 'rel/tree'
    }).then((res) => {
      if (res.data == null) {
        state.tree = []
      } else {
        state.tree = res.data
      }
      state.loading = false
      state.loaded = true
    }).catch((err) => {
      Message.error(err)
      this.loading = false
    })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
