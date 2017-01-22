import Vue from 'vue'
import VueRouter from 'vue-router'
import login from './components/login'
import ldap from './components/login/ldap'
import misso from './components/login/misso'
import meta from './components/meta'
import expr from './components/meta/expression'
import exprEdit from './components/meta/expression/edit'
import host from './components/meta/host'
import hostEdit from './components/meta/host/edit'
import role from './components/meta/role'
import roleEdit from './components/meta/role/edit'
import rule from './components/meta/rule'
import ruleEdit from './components/meta/rule/edit'
import tag from './components/meta/tag'
import tagEdit from './components/meta/tag/edit'
import team from './components/meta/team'
import teamEdit from './components/meta/team/edit'
import token from './components/meta/token'
import tokenEdit from './components/meta/token/edit'
import user from './components/meta/user'
import userEdit from './components/meta/user/edit'
import store from 'store'
const { _ } = window

Vue.use(VueRouter)

const router = new VueRouter({
  routes:
  [{
    path: '/login',
    redirect: '/login/ldap',
    component: login,
    children: [
    { path: 'ldap', component: ldap },
    { path: 'misso', component: misso }
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
    { path: 'rule', redirect: 'rule/list' },
    { path: 'rule/list', component: rule },
    { path: 'rule/edit', component: ruleEdit },
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
  }]
})

router.beforeEach((to, from, next) => {
  if (store.state.login.status) {
    console.log('1', store.state.login.status, 'next')
    next()
    return
  }

  if (_.startsWith(to.path, '/login') ||
    _.startsWith(to.path, '/help')) {
    console.log('2', store.state.login.status, 'next')
    next()
    return
  }

  console.log('before call login', store.state.login.status)
  store.dispatch('login_quiet')
  next()
  setTimeout(() => {
    if (!store.state.login.status) {
      console.log('4', store.state.login.status, 'next')
      next('/login/ldap')
    }
  }, 2000)
})

export default router
