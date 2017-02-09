import Vue from 'vue'
import VueRouter from 'vue-router'
import login from './components/login'
import ldap from './components/login/ldap'
import misso from './components/login/misso'
import settings from './components/settings'
import ctrl from './components/settings/ctrl'
import agent from './components/settings/agent'
import lb from './components/settings/lb'
import backend from './components/settings/backend'
import profile from './components/settings/profile'
import aboutme from './components/settings/aboutme'
import debug from './components/settings/debug'
import meta from './components/meta'
import expr from './components/meta/expression'
import exprEdit from './components/meta/expression/edit'
import host from './components/meta/host'
import hostEdit from './components/meta/host/edit'
import role from './components/meta/role'
import roleEdit from './components/meta/role/edit'
import template from './components/meta/template'
import templateEdit from './components/meta/template/edit'
import tag from './components/meta/tag'
import tagEdit from './components/meta/tag/edit'
import team from './components/meta/team'
import teamEdit from './components/meta/team/edit'
import token from './components/meta/token'
import tokenEdit from './components/meta/token/edit'
import user from './components/meta/user'
import userEdit from './components/meta/user/edit'
import rel from './components/rel'
import tagHost from './components/rel/tag_host'
import tagRoleUser from './components/rel/tag_role_user'
import tagRoleToken from './components/rel/tag_role_token'
import tagTemplateTrigger from './components/rel/tag_template_trigger'
import store from './store'
const { _ } = window

Vue.use(VueRouter)

const router = new VueRouter({
  routes:
  [{
    path: '/',
    redirect: '/meta/tag/list'
  }, {
    path: '/login',
    redirect: '/login/ldap',
    component: login,
    children: [
    { path: 'ldap', component: ldap },
    { path: 'misso', component: misso }
    ]
  }, {
    path: '/settings',
    redirect: '/settings/aboutme',
    component: settings,
    children: [
    { path: 'config/ctrl', component: ctrl },
    { path: 'config/agent', component: agent },
    { path: 'config/lb', component: lb },
    { path: 'config/backend', component: backend },
    { path: 'profile', component: profile },
    { path: 'aboutme', component: aboutme },
    { path: 'debug', component: debug }
    ]
  }, {
    path: '/meta',
    redirect: '/meta/host/list',
    component: meta,
    children: [
    { path: 'expression', redirect: 'expression/list' },
    { path: 'expression/list', component: expr },
    { path: 'expression/edit', component: exprEdit },
    { path: 'host', redirect: 'host/list' },
    { path: 'host/list', component: host },
    { path: 'host/edit', component: hostEdit },
    { path: 'role', redirect: 'role/list' },
    { path: 'role/list', component: role },
    { path: 'role/edit', component: roleEdit },
    { path: 'template', redirect: 'template/list' },
    { path: 'template/list', component: template },
    { path: 'template/edit', component: templateEdit },
    { path: 'tag', redirect: 'tag/list' },
    { path: 'tag/list', component: tag },
    { path: 'tag/edit', component: tagEdit },
    { path: 'team', redirect: 'team/list' },
    { path: 'team/list', component: team },
    { path: 'team/edit', component: teamEdit },
    { path: 'token', redirect: 'token/list' },
    { path: 'token/list', component: token },
    { path: 'token/edit', component: tokenEdit },
    { path: 'user', redirect: 'user/list' },
    { path: 'user/list', component: user },
    { path: 'user/edit', component: userEdit }
    ]
  }, {
    path: '/rel',
    redirect: '/rel/tag-host',
    component: rel,
    children: [
    { path: 'tag-host', component: tagHost },
    { path: 'tag-role-user', component: tagRoleUser },
    { path: 'tag-role-token', component: tagRoleToken },
    { path: 'tag-template-trigger', component: tagTemplateTrigger }
    ]
  }]
})

router.beforeEach((to, from, next) => {
  console.log(to)
  if (store.state.login.login) {
    next()
    return
  }

  if (_.startsWith(to.path, '/login') ||
    _.startsWith(to.path, '/help')) {
    next()
    return
  }

  if (window.Cookies.get('username') !== undefined) {
    next(vm => { vm.$store.dispatch('login/login', {router: vm.$router}) })
    return
  }

  // next('/login/ldap')
  next(vm => {
    vm.$store.commit('login/m_set_callback', vm.$router.fullPath)
    vm.$router.push({path: '/login'})
  })
})

export default router
