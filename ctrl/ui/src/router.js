// resources:
// https://github.com/vuejs/vue-router/blob/next/examples/redirect/app.js
import Vue from 'vue'
import VueRouter from 'vue-router'

import login from './components/login'

import admin from './components/admin'
import ctrl from './components/admin/ctrl'
import agent from './components/admin/agent'
import loadbalance from './components/admin/loadbalance'
import backend from './components/admin/backend'
import debug from './components/admin/debug'

import settings from './components/settings'
import profile from './components/settings/profile'
import about from './components/settings/about'
import log from './components/settings/log'

import meta from './components/meta'
import expression from './components/meta/expression'
import host from './components/meta/host'
import role from './components/meta/role'
import tag from './components/meta/tag'
import team from './components/meta/team'
import template from './components/meta/template'
import token from './components/meta/token'
import user from './components/meta/user'

import relation from './components/relation'
import tagHost from './components/relation/tag_host'
import tagRoleUser from './components/relation/tag_role_user'
import tagRoleToken from './components/relation/tag_role_token'
import tagTemplate from './components/relation/tag_template'

Vue.use(VueRouter)

import store from 'store'
import { Msg } from 'src/utils'

var got = false

function access (type, to, from, next) {
  if (!got) {
    got = true
    store.dispatch('auth/info').then(() => {
      _access(type, to, from, next)
      if (store.state.auth.login) {
        store.dispatch('load_config')
      }
    }).catch(() => {
      Msg.error('not login')
      next({path: '/login', query: {cb: to.fullPath}})
    })
    return
  }

  _access(type, to, from, next)
}

function _access (type, to, from, next) {
  let x = false

  if (type === 'login') {
    x = store.state.auth.login
  } else if (type === 'reader') {
    x = store.state.auth.reader
  } else if (type === 'admin') {
    x = store.state.auth.admin
  }

  if (x) {
    next()
  } else {
    Msg.error('not ' + type)
    if (!store.state.auth.login) {
      next({path: '/login', query: {cb: to.fullPath}})
    } else {
      next(false)
    }
  }
}

function accessLogin (to, from, next) {
  access('login', to, from, next)
}

function accessReader (to, from, next) {
  access('reader', to, from, next)
}

function accessAdmin (to, from, next) {
  access('admin', to, from, next)
}

const router = new VueRouter({
  routes:
  [{
    path: '/',
    beforeEnter: accessLogin,
    redirect: '/settings'
  }, {
    path: '/login',
    component: login
  }, {
    path: '/admin',
    redirect: '/admin/ctrl',
    component: admin,
    beforeEnter: accessAdmin,
    children: [
    { path: 'ctrl', component: ctrl },
    { path: 'agent', component: agent },
    { path: 'loadbalance', component: loadbalance },
    { path: 'backend', component: backend },
    { path: 'profile', component: profile },
    { path: 'debug', component: debug }
    ]
  }, {
    path: '/settings',
    redirect: '/settings/about',
    component: settings,
    beforeEnter: accessLogin,
    children: [
    { path: 'profile', component: profile },
    { path: 'about', component: about },
    { path: 'log', component: log }
    ]
  }, {
    path: '/meta',
    redirect: '/meta/host',
    component: meta,
    beforeEnter: accessReader,
    children: [
    { path: 'expression', component: expression },
    { path: 'host', component: host },
    { path: 'role', component: role },
    { path: 'tag', component: tag },
    { path: 'team', component: team },
    { path: 'template', component: template },
    { path: 'token', component: token },
    { path: 'user', component: user }
    ]
  }, {
    path: '/relation',
    redirect: '/relation/tag-host',
    component: relation,
    beforeEnter: accessReader,
    children: [
    { path: 'tag-host', component: tagHost },
    { path: 'tag-role-user', component: tagRoleUser },
    { path: 'tag-role-token', component: tagRoleToken },
    { path: 'tag-template', component: tagTemplate }
    ]
  }, { path: '*', redirect: '/' }
  ]
})

export default router
