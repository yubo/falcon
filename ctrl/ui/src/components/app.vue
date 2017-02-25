<template>
<div id="app">
  <!-- Fixed navbar -->
  <nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container">
      <div class="navbar-header"> <a class="navbar-brand" href="#">Falcon</a> </div>
      <div id="navbar" class="navbar-collapse collapse">
        <ul class="nav navbar-nav">
          <li is="li-tpl" v-for="(nav, li_idx) in navs" :obj="nav"></li>
        </ul>
        <ul class="nav navbar-nav navbar-right">
          <li class="dropdown" v-if="login">
            <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
              aria-expanded="false">{{this.$store.state.auth.username}}<span class="caret"></span></a>
            <ul class="dropdown-menu">
              <li><router-link to="/settings/profile">Profile</router-link></li>
              <li><router-link to="/settings/about">About</router-link></li>
              <li><a href="/doc" target="_blank">doc</a></li>
              <li><a href="#" @click="logout">logout</a></li>
            </ul>
          </li>
        </ul>
      </div><!--/.nav-collapse -->
    </div>
  </nav>
  <div class="container">
    <div class="container theme-showcase" role="main" v-if="login">
      <ul class="nav nav-tabs" role="tablist">
        <li is="li-tpl" v-for="(nav, li_idx) in subnavs" :obj="nav"></li>
      </ul>
    </div>
  
    <router-view></router-view>
  </div>

</div>
</template>

<script>
import { liTpl } from './tpl'

export default {
  components: {
    liTpl
  },
  data () {
    return {
      links: {
        dashboard: [
        {url: '/dashboard/graph', text: 'grpah'},
        {url: '/dashboard/alarm', text: 'alarm'}
        ],
        relation: [
        {url: '/relation/tag-host', text: 'Tag Host'},
        {url: '/relation/tag-template', text: 'Tag Template'},
        {url: '/relation/tag-role-user', text: 'Tag Role User'},
        {url: '/relation/tag-role-token', text: 'Tag Role Token'}
        ],
        meta: [
        {url: '/meta/tag', text: 'Tag'},
        {url: '/meta/host', text: 'Host'},
        {url: '/meta/role', text: 'Role'},
        {url: '/meta/user', text: 'User'},
        {url: '/meta/token', text: 'token'},
        {url: '/meta/team', text: 'Team'},
        {url: '/meta/template', text: 'Template'},
        {url: '/meta/expression', text: 'Expression'}
        ],
        admin: [
        {url: '/admin/ctrl', text: 'Ctrl'},
        {url: '/admin/agent', text: 'Agent'},
        {url: '/admin/loadbalance', text: 'LoadBalance'},
        {url: '/admin/backend', text: 'Backend'},
        {url: '/admin/debug', text: 'Debug'}
        ],
        settings: [
        {url: '/settings/about', text: 'About'},
        {url: '/settings/profile', text: 'Profile'},
        {url: '/settings/log', text: 'Log'}
        ]
      }
    }
  },
  computed: {
    login () {
      return this.$store.state.auth.login
    },
    isAdmin () {
      return this.$store.state.auth.admin
    },
    isReader () {
      return this.$store.state.auth.reader
    },
    isOperator () {
      return this.$store.state.auth.operator
    },
    navs () {
      let links = []
      if (this.isReader) {
        // links.push({url: '/dashboard', text: 'Dashboard'})
        links.push({url: '/relation', text: 'Relatioin'})
        links.push({url: '/meta', text: 'Meta'})
      }
      if (this.isAdmin) {
        links.push({url: '/admin', text: 'Admin'})
      }
      if (this.login) {
        links.push({url: '/settings', text: 'Settings'})
      }
      return links
    },
    subnavskey () {
      return this.$route.path.slice(1).split('/')[0]
    },
    subnavs () {
      return this.links[this.$route.path.slice(1).split('/')[0]]
      // return []
    }
  },
  created () {
  },
  methods: {
    logout () {
      this.$store.dispatch('auth/logout',
            {router: this.$router, cb: '/'})
    }
  }
}
</script>

<style>
.nav-tabs {
margin:10px;
}
</style>
