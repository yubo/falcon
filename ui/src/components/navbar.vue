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

        <li class="dropdown" v-if="login">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Dashboard<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/dashboard/falcon">Falcon</router-link></li>
            <li><router-link to="/dashborad/alarm">Alarm</router-link></li>
            <li><router-link to="/dashborad/graph">Graph</router-link></li>
          </ul>
        </li>
        <li class="dropdown" v-if="login">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Relation<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/rel/tag-host">Tag Host</router-link></li>
            <li><router-link to="/rel/tag-role-user">Tag Role User</router-link></li>
            <li><router-link to="/rel/tag-role-token">Tag Role Token</router-link></li>
            <li><router-link to="/rel/tag-rule-trigger">Tag Template Trigger</router-link></li>
          </ul>
        </li>
        <li class="dropdown">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Meta<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/meta/tag/list">Tag</router-link></li>
            <li><router-link to="/meta/host/list">Host</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/meta/role/list">Role</router-link></li>
            <li><router-link to="/meta/user/list">User</router-link></li>
            <li><router-link to="/meta/token/list">Token</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/meta/team/list">Team</router-link></li>
            <li><router-link to="/meta/rule/list">Template</router-link></li>
            <li class="disabled"><router-link to="#">Trigger</router-link></li>
            <li><router-link to="/meta/expression/list">Expression</router-link></li>
          </ul>
        </li>
        <li class="dropdown" v-if="login">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Settings<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/settings/config/ctrl">ctrl</router-link></li>
            <li><router-link to="/settings/config/agent">agent</router-link></li>
            <li><router-link to="/settings/config/graph">graph</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/settings/profile">Profile</router-link></li>
            <li><router-link to="/settings/aboutme">About Me</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/settings/debug">Debug</router-link></li>
          </ul>
        </li>
        <li class="dropdown">
          <a to="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">help<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><router-link to="/doc">doc</router-link></li>
            <li role="separator" class="divider"></li>
            <li><router-link to="/about">About Falcon</router-link></li>
        </li>
        <li v-if="login"><a @click="logout">[logout]</a></li>
        <li v-if="!login"><router-link to="/login/ldap">[login]</router-link></li>
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
      return this.$store.state.login.status
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
      this.$store.dispatch('logout')
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
