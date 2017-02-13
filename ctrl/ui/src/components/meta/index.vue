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
      { url: '/meta/tag', text: 'Tag' },
      { url: '/meta/host', text: 'Host' }
      ], [
      { url: '/meta/role', text: 'Role' },
      { url: '/meta/user', text: 'User' },
      { url: '/meta/token', text: 'Token' }
      ], [
      { url: '/meta/team', text: 'Team' },
      { url: '/meta/template', text: 'Template' },
      { url: '/meta/expression', text: 'Expression' }
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
