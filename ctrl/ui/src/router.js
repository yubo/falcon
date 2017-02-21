import Vue from 'vue'
import VueRouter from 'vue-router'

import admin from './components/admin'
import ctrl from './components/admin/ctrl'
import agent from './components/admin/agent'
import loadbalance from './components/admin/loadbalance'
import backend from './components/admin/backend'
import debug from './components/admin/debug'

import settings from './components/settings'
import profile from './components/settings/profile'
import about from './components/settings/about'

import meta from './components/meta'
import expression from './components/meta/expression'
import host from './components/meta/host'
import role from './components/meta/role'
import tag from './components/meta/tag'
import team from './components/meta/team'
import template from './components/meta/template'
import token from './components/meta/token'
import user from './components/meta/user'

import rel from './components/rel'
import tagHost from './components/rel/tag_host'
import tagRoleUser from './components/rel/tag_role_user'
import tagRoleToken from './components/rel/tag_role_token'
import tagTemplate from './components/rel/tag_template'

Vue.use(VueRouter)

import store from 'store'
import { Msg } from 'src/utils'

var got = false

function gotInfo () {
  got = true
  store.commit('auth/m_set_loading', true)
  store.dispatch('auth/info')
}

function accessReader (to, from, next) {
  if (!got) {
    gotInfo()
  }
  if (!store.state.auth.loading) {
    if (store.state.auth.reader) {
      next()
    } else {
      Msg.error('permission denied')
      next(false)
    }
    return
  }

  setTimeout(() => {
    accessReader(to, from, next)
  }, 100)
}

function accessAdmin (to, from, next) {
  if (!got) {
    gotInfo()
  }
  if (!store.state.auth.loading) {
    if (store.state.auth.admin) {
      next()
    } else {
      Msg.error('permission denied')
      next(false)
    }
    return
  }

  setTimeout(() => {
    accessAdmin(to, from, next)
  }, 100)
}

const router = new VueRouter({
  routes:
  [{
    path: '/',
    redirect: '/settings/about'
  }, {
    path: '/admin',
    redirect: '/admin/ctrl',
    component: admin,
    beforeEnter: accessAdmin,
    children: [
    { path: 'config/ctrl', component: ctrl },
    { path: 'config/agent', component: agent },
    { path: 'config/loadbalance', component: loadbalance },
    { path: 'config/backend', component: backend },
    { path: 'profile', component: profile },
    { path: 'debug', component: debug }
    ]
  }, {
    path: '/settings',
    redirect: '/settings/about',
    component: settings,
    beforeEnter: accessReader,
    children: [
    { path: 'profile', component: profile },
    { path: 'about', component: about }
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
    path: '/rel',
    redirect: '/rel/tag-host',
    component: rel,
    beforeEnter: accessReader,
    children: [
    { path: 'tag-host', component: tagHost },
    { path: 'tag-role-user', component: tagRoleUser },
    { path: 'tag-role-token', component: tagRoleToken },
    { path: 'tag-template', component: tagTemplate }
    ]
  }]
})

export default router
