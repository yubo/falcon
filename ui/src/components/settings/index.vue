<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-sm-3 col-md-2 sidebar">
        <ul v-for="(objs, ul_idx) in links" class="nav nav-sidebar">
          <li is="li-tpl" v-for="(obj, li_idx) in objs" :obj="obj"></li>
        </ul>
      </div>
      <router-view> </router-view>
  </div>
</template>

<script>
import { liTpl } from '../tpl'

export default {
  data () {
    return {
      links:
      [[
      { url: '/settings/config/ctrl', text: 'Config Ctrl' },
      { url: '/settings/config/agent', text: 'Config Agent' },
      { url: '/settings/config/lb', text: 'Config Lb' },
      { url: '/settings/config/backend', text: 'Config Backend' }
      ], [
      { url: '/settings/profile', text: 'Profile' },
      { url: '/settings/aboutme', text: 'About Me' }
      ], [
      { url: '/settings/debug', text: 'Debug' }
      ]]
    }
  },
  components: {
    liTpl
  },
  computed: {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
