<template>
  <div class="container-fluid" v-if="isAdmin">
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
      { url: '/admin/config/ctrl', text: 'Ctrl' },
      { url: '/admin/config/agent', text: 'Agent' },
      { url: '/admin/config/loadbalance', text: 'Load Balance' },
      { url: '/admin/config/backend', text: 'Backend' }
      ], [
      { url: '/admin/debug', text: 'Debug' }
      ]]
    }
  },
  components: {
    liTpl
  },
  computed: {
    isAdmin () {
      return this.$store.state.auth.admin
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
