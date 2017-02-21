<template>
<nav class="navbar navbar-inverse navbar-fixed-top">
  <div class="container-fluid">
    <div class="navbar-header">
      <ol class="breadcrumb">
         <li is="li-tpl" v-for="(obj, index) in links" :obj="obj"></li>
      </ol>
    </div>
    <div id="navbar" class="navbar-collapse collapse">
      <ul class="nav navbar-nav navbar-right">

        <li class="dropdown" v-if="isReader">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Dashboard<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/dashboard/falcon">Falcon</router-link></li>
            <li><router-link to="/dashborad/alarm">Alarm</router-link></li>
            <li><router-link to="/dashborad/graph">Graph</router-link></li>
          </ul>
        </li>
        <li class="dropdown" v-if="isReader">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Relation<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/rel/tag-host">Tag Host</router-link></li>
            <li><router-link to="/rel/tag-template">Tag Template</router-link></li>
            <li><router-link to="/rel/tag-role-user">Tag Role User</router-link></li>
            <li><router-link to="/rel/tag-role-token">Tag Role Token</router-link></li>
          </ul>
        </li>
        <li class="dropdown" v-if="isReader">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Meta<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/meta/tag">Tag</router-link></li>
            <li><router-link to="/meta/host">Host</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/meta/role">Role</router-link></li>
            <li><router-link to="/meta/user">User</router-link></li>
            <li><router-link to="/meta/token">Token</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/meta/team">Team</router-link></li>
            <li><router-link to="/meta/template">Template</router-link></li>
            <li><router-link to="/meta/expression">Expression</router-link></li>
          </ul>
        </li>
        <li class="dropdown" v-if="isAdmin">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Admin<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/admin/config/ctrl">Ctrl</router-link></li>
            <li><router-link to="/admin/config/agent">Agent</router-link></li>
            <li><router-link to="/admin/config/loadbalance">Load Balance</router-link></li>
            <li><router-link to="/admin/config/backend">Backend</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/admin/debug">Debug</router-link></li>
          </ul>
        </li>
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
    </div>
  </div>
</nav>
</template>

<script>
var liTpl = {
  template: `<li v-if="obj.text !== ''"> <router-link :to="obj.url">{{ obj.text }}</router-link></li>`,
  props: ['obj']
}

export default {
  components: {
    liTpl
  },
  data () {
    return {
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
    links () {
      var url = ''
      return this.$route.path.slice(1).split('/').map((link) => {
        url = url + '/' + link
        return {
          url: url,
          text: link.charAt(0).toUpperCase() + link.slice(1)
        }
      })
    }
  },
  methods: {
    logout () {
      this.$store.dispatch('auth/logout')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
.navbar-header > .breadcrumb {
  background-color: transparent;
  margin:0px 0px 0px 0px;
  padding: 15px 15px 15px 15px;
}
.navbar-header > .breadcrumb li a{
   font-size:14px;
   color: #9d9d9d;
   height:20px;
}
</style>
