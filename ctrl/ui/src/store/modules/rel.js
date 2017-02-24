import { fetch } from 'src/utils'
import { Message } from 'element-ui'
const { _ } = window

const state = {
  loading: false,
  loaded: false,
  curTag: {name: '', id: 1},
  tree: [],
  opnode: []
}

const getters = {
}

const actions = {
  load_opnode ({state, commit}) {
    return new Promise((resolve, reject) => {
      fetch({
        method: 'get',
        url: 'rel/tree/opnode'
      }).then((res) => {
        resolve(res)
      }).catch((err) => {
        reject(err)
      })
    })
  },
  load_tree ({state, commit, rootState, dispatch}) {
    commit('m_set_loading', true)
    fetch({
      method: 'get',
      url: 'rel/tree'
    }).then((res) => {
      if (rootState.auth.admin) {
        commit('m_set_tree', {
          tree: res.data,
          admin: rootState.auth.admin,
          opnode: {}
        })
        commit('m_set_loading', false)
      } else {
        dispatch('load_opnode').then((res2) => {
          commit('m_set_tree', {
            tree: res.data,
            admin: rootState.auth.admin,
            opnode: res2.data
          })
          commit('m_set_loaded', true)
          commit('m_set_loading', false)
        }).catch((err) => {
          Message.error(err)
          commit('m_set_loading', false)
        })
      }
    }).catch((err) => {
      Message.error(err)
      commit('m_set_loading', false)
    })
  }
}

function pruneTree (tree, opnode) {
  if (!tree) {
    return
  }
  for (let n in tree) {
    let array = []
    pruneTree(tree[n].child, opnode)
    tree[n].child = _.remove(array, (o) => {
      return ((!o.child || o.child.length === 0) && !opnode[o.id])
    })
    console.log(tree[n].name, 'remove', array)
  }
}

const mutations = {
  'm_cur_tag' (state, val) {
    state.curTag = val
  },
  'm_set_loading' (state, val) {
    state.loading = val
  },
  'm_set_loaded' (state, val) {
    state.loaded = val
  },
  'm_set_tree' (state, args) {
    state.tree = args.tree
    console.log(args)
    if (!args.admin) {
      let node = {}
      for (let k in args.opnode) {
        node[args.opnode[k]] = true
      }
      pruneTree(state.tree, node)
      console.log('tree', state.tree)
    }
  },
  'm_set_opnode' (state, val) {
    state.opnode = val
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
