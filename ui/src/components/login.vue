<template>
<div class="container">
  <ul class="nav nav-tabs">
    <li role="presentation" v-bind:class="{ active: tabName === 'ldap' }"><a @click="showLdap()" >ldap</a></li>
    <li role="presentation" v-bind:class="{ active: tabName === 'misso' }"><a @click="showMisso()" >misso</a></li>
  </ul>
  <router-view> </router-view> 
</div>
</template>

<script>
import ldap from './login/ldap'
import misso from './login/misso'

export default {
  data () {
    return {
      tabName: 'ldap'
    }
  },
  components: {
    ldap,
    misso
  },
  methods: {
    showLdap () {
      this.tabName = 'ldap'
      this.$router.push('/login/ldap')
    },
    showMisso () {
      this.tabName = 'misso'
      this.$router.push('/login/misso')
    }
  }

}
</script>

<style>
</style>
