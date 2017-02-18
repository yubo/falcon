import Vue from 'vue'
import VueRouter from 'vue-router'

import settings from './components/settings'
import ctrl from './components/settings/ctrl'
import agent from './components/settings/agent'
import loadbalance from './components/settings/loadbalance'
import backend from './components/settings/backend'
import profile from './components/settings/profile'
import about from './components/settings/about'
import debug from './components/settings/debug'

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

const router = new VueRouter({
  routes:
  [{
    path: '/',
    redirect: '/meta/tag'
  }, {
    path: '/settings',
    redirect: '/settings/about',
    component: settings,
    children: [
    { path: 'config/ctrl', component: ctrl },
    { path: 'config/agent', component: agent },
    { path: 'config/loadbalance', component: loadbalance },
    { path: 'config/backend', component: backend },
    { path: 'profile', component: profile },
    { path: 'about', component: about },
    { path: 'debug', component: debug }
    ]
  }, {
    path: '/meta',
    redirect: '/meta/host',
    component: meta,
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
    children: [
    { path: 'tag-host', component: tagHost },
    { path: 'tag-role-user', component: tagRoleUser },
    { path: 'tag-role-token', component: tagRoleToken },
    { path: 'tag-template', component: tagTemplate }
    ]
  }]
})

export default router
