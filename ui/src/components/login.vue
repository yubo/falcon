<template>
<div class="container">
  <ul class="nav nav-tabs">
    <li role="presentation" v-bind:class="{ active: tabName === 'ldap' }"><a @click="tabName = 'ldap'" >ldap</a></li>
    <li role="presentation" v-bind:class="{ active: tabName === 'misso' }"><a @click="tabName = 'misso'" >misso</a></li>
  </ul>
  <ldap v-show="tabName === 'ldap'"> </ldap> 
  <misso v-show="tabName === 'misso'"></misso>
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
  }

}
</script>

<style>
</style>
