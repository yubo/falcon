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
    pruneTree(tree[n].child, opnode)
    _.remove(tree[n].child, (o) => {
      return ((!o.child || o.child.length === 0) && !opnode[o.id])
    })
    if (!opnode[tree[n].id]) {
      tree[n].label += ' *'
      tree[n].ro = true
    }
  }
}

function getNode (node, id) {
  if (node.id === id) {
    return node
  }
  for (let n in node.child) {
    const o = getNode(node.child[n], id)
    if (o) {
      return o
    }
  }
  return null
}

function delNode (node, id, g) {
  if (g.done) {
    return
  }
  let a = _.remove(node.child, (n) => {
    return n.id === id
  })
  if (a.length > 0) {
    g.done = true
    return
  }
  for (let n in node.child) {
    delNode(node.child[n], id, g)
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
    if (!args.admin) {
      let node = {}
      for (let k in args.opnode) {
        node[args.opnode[k]] = true
      }
      pruneTree(state.tree, node)
    }
  },
  'm_add_node' (state, args) {
    let node = getNode(state.tree[0], args.id)
    if (!node.child) {
      node.child = []
    }
    node.child.push({
      id: args.cid,
      label: args.label,
      name: args.name,
      child: null,
      ro: false
    })
  },
  'm_del_node' (state, id) {
    let g = {done: false}
    delNode(state.tree[0], id, g)
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
