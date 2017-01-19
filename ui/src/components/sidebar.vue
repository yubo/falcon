<template>
<div v-if="tree" class="col-sm-3 col-md-2 sidebar">
  <div class="zTreeDemoBackground left">
    <ul id="tree" class="ztree"></ul>
    hello
  </div>
</div>
<div v-else class="col-sm-3 col-md-2 sidebar">
  <ul class="nav nav-sidebar">
    <li><a href="/settings/config/ctrl">ctrl {{ msg1 }}</a> </li>
    <li><a href="/settings/config/agent">agent</a> </li>
    <li><a href="/settings/config/graph">graph</a> </li>
  </ul>
  <ul class="nav nav-sidebar">
    <li><a href="/settings/profile">Profile</a> </li>
    <li class="active"><a href="#">About Me<span class="sr-only">(current)</span></a></li>
  </ul>
  <ul class="nav nav-sidebar">
    <li><a href="/settings/debug">Debug</a> </li>
  </ul>
</template>

<script>
import { mapState } from 'vuex'

export default {
  computed: mapState({
    tree: state => state.tree.show
  }),
  data () {
    return {
      bbb: false,
      msg1: 'hello1',
      msg: 'Welcome to falcon@xiaomi.com App'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
